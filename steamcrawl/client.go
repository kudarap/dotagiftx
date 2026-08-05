// Package steamcrawl provides a Steam-aware HTTP client that evades
// rate limit (429) errors using browser-like sessions, exponential
// backoff with jitter and honor of server rate limit headers.
package steamcrawl

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// Config configures the crawler client behavior.
type Config struct {
	// MaxAttempts is the total number of request attempts including the
	// initial one. Defaults to 5.
	MaxAttempts int
	// BaseBackoff is the initial backoff duration. Defaults to 1s.
	BaseBackoff time.Duration
	// MaxBackoff caps the backoff duration. Defaults to 30s.
	MaxBackoff time.Duration
	// Timeout is the per attempt http request timeout. Defaults to 15s.
	Timeout time.Duration
	// RetryOnStatus lists http status codes that trigger a retry.
	// Defaults to 429, 500, 502, 503 and 504.
	RetryOnStatus []int
	// HTTPClient overrides the underlying http client. When nil a client
	// with the configured Timeout is used.
	HTTPClient *http.Client
}

// DefaultConfig returns a config with sane defaults.
func DefaultConfig() Config {
	return Config{
		MaxAttempts:   5,
		BaseBackoff:   time.Second,
		MaxBackoff:    30 * time.Second,
		Timeout:       15 * time.Second,
		RetryOnStatus: []int{http.StatusTooManyRequests, 500, 502, 503, 504},
	}
}

// Client crawls steam endpoints while evading rate limit errors.
type Client struct {
	cfg  Config
	http *http.Client
}

// New returns a crawler client from config applying defaults.
func New(cfg Config) *Client {
	cfg = cfg.setDefault()

	hc := cfg.HTTPClient
	if hc == nil {
		hc = &http.Client{Timeout: cfg.Timeout}
	}
	return &Client{cfg: cfg, http: hc}
}

// Do performs the request retrying on rate limit and transient errors.
// On success it returns the response status code and body. When out is
// non nil and the body is valid JSON it unmarshals into it. It returns
// ErrTooManyRequests when all attempts were rate limited.
func (c *Client) Do(ctx context.Context, req *http.Request, out any) (int, []byte, error) {
	var lastStatus int
	var lastErr error

	for attempt := 1; attempt <= c.cfg.MaxAttempts; attempt++ {
		status, body, retryAfter, err := c.doOnce(ctx, req)
		lastStatus, lastErr = status, err

		if err == nil {
			if out != nil && len(body) > 0 {
				if err = json.Unmarshal(body, out); err != nil {
					return status, body, err
				}
			}
			return status, body, nil
		}

		// stop when the status is not meant to be retried.
		if status != 0 && !c.isRetryable(status) {
			return status, body, err
		}
		if attempt == c.cfg.MaxAttempts {
			break
		}

		delay := computeDelay(attempt, retryAfter, c.cfg.BaseBackoff, c.cfg.MaxBackoff)
		select {
		case <-ctx.Done():
			return status, body, ctx.Err()
		case <-time.After(delay):
		}
	}

	if lastStatus == http.StatusTooManyRequests {
		return lastStatus, nil, ErrTooManyRequests
	}
	if lastErr == nil {
		lastErr = ErrRetriesExhausted
	}
	return lastStatus, nil, lastErr
}

// doOnce sends a single attempt with a fresh session. It returns the
// retryAfter delay when the server asks us to wait.
func (c *Client) doOnce(ctx context.Context, req *http.Request) (status int, body []byte, retryAfter time.Duration, err error) {
	r := req.Clone(ctx)
	c.applySession(r)

	res, err := c.http.Do(r)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return 0, nil, 0, ctxErr
		}
		return 0, nil, 0, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, nil, 0, err
	}

	if res.StatusCode >= 400 {
		retryAfter = retryAfterHeader(res.Header)
		return res.StatusCode, b, retryAfter, &statusError{res.StatusCode, res.Status}
	}

	return res.StatusCode, b, 0, nil
}

// applySession attaches a fresh browser-like session to the request.
func (c *Client) applySession(r *http.Request) {
	steamID := SteamIDFrom(r.Context())
	if steamID == "" {
		steamID = r.URL.Query().Get("steam_id")
	}
	session := NewSession(steamID)
	for k, vs := range session {
		for _, v := range vs {
			if k == "Cookie" {
				r.Header.Add("Cookie", v)
				continue
			}
			r.Header.Set(k, v)
		}
	}
}

func (c *Client) isRetryable(status int) bool {
	for _, s := range c.cfg.RetryOnStatus {
		if status == s {
			return true
		}
	}
	return false
}

func (c Config) setDefault() Config {
	if c.MaxAttempts <= 0 {
		c.MaxAttempts = 5
	}
	if c.BaseBackoff <= 0 {
		c.BaseBackoff = time.Second
	}
	if c.MaxBackoff <= 0 {
		c.MaxBackoff = 30 * time.Second
	}
	if c.Timeout <= 0 {
		c.Timeout = 15 * time.Second
	}
	if len(c.RetryOnStatus) == 0 {
		c.RetryOnStatus = []int{http.StatusTooManyRequests, 500, 502, 503, 504}
	}
	return c
}
