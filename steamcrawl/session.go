package steamcrawl

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strings"
	"time"
)

const (
	uaChrome  = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
	uaFirefox = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) Gecko/20100101 Firefox/125.0"
	uaSafari  = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4 Safari/605.1.15"
	uaEdge    = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36 Edg/124.0.0.0"

	uaLinuxChrome = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
	uaLinuxFox    = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:125.0) Gecko/20100101 Firefox/125.0"
)

var userAgents = []string{uaChrome, uaFirefox, uaSafari, uaEdge, uaLinuxChrome, uaLinuxFox}

var acceptLanguages = []string{
	"en-US,en;q=0.9",
	"en-GB,en;q=0.8",
	"en-CA,en;q=0.9",
	"en-AU,en;q=0.9",
	"en,en-US;q=0.9",
}

// NewSession builds a browser-like header set for a steam community
// request. Each call returns a fresh session with randomized user agent,
// accept language and session cookies so retries look like a new browser.
func NewSession(steamID string) http.Header {
	h := http.Header{}
	h.Set("User-Agent", randomChoice(userAgents))
	h.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	h.Set("Accept-Language", randomChoice(acceptLanguages))
	h.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	h.Set("Connection", "keep-alive")
	h.Set("Referer", referer(steamID))
	h.Set("Sec-Fetch-Dest", "empty")
	h.Set("Sec-Fetch-Mode", "cors")
	h.Set("Sec-Fetch-Site", "same-origin")
	h.Set("Sec-GPC", "1")
	h.Set("Pragma", "no-cache")
	h.Set("Cache-Control", "no-cache")

	h.Add("Cookie", freshSessionCookie())
	return h
}

func referer(steamID string) string {
	if steamID == "" {
		return "https://steamcommunity.com/"
	}
	return fmt.Sprintf("https://steamcommunity.com/profiles/%s/inventory/", steamID)
}

// WithSteamID attaches a steam id to the context so the session referer
// points to the correct profile inventory.
func WithSteamID(ctx context.Context, steamID string) context.Context {
	return context.WithValue(ctx, steamIDCtxKey{}, steamID)
}

// SteamIDFrom returns the steam id attached by WithSteamID.
func SteamIDFrom(ctx context.Context) string {
	v, _ := ctx.Value(steamIDCtxKey{}).(string)
	return v
}

type steamIDCtxKey struct{}

// freshSessionCookie returns a randomized session id cookie so each
// retry attempt is treated as a brand new visitor session.
func freshSessionCookie() string {
	return fmt.Sprintf("sessionid=%s; SteamLanguage=english", randString(24))
}

func randomChoice[T any](items []T) T {
	return items[rand.IntN(len(items))]
}

func randString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(letters[rand.IntN(len(letters))])
	}
	return sb.String()
}

// retryAfterHeader extracts a delay from common rate limit response
// headers.
func retryAfterHeader(h http.Header) time.Duration {
	if v := h.Get("Retry-After"); v != "" {
		return parseRetryAfter(v)
	}
	if v := h.Get("RateLimit-Reset"); v != "" {
		return parseRetryAfter(v)
	}
	return 0
}
