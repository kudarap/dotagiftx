package steamcrawl

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestDo_retriesOnTooManyRequests(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) <= 3 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	c := New(Config{MaxAttempts: 5, BaseBackoff: time.Millisecond, MaxBackoff: time.Millisecond})
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
	var out struct {
		Ok bool `json:"ok"`
	}
	status, body, err := c.Do(context.Background(), req, &out)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("Do() status = %d, want 200", status)
	}
	if calls.Load() != 4 {
		t.Fatalf("Do() calls = %d, want 4", calls.Load())
	}
	if !out.Ok || !strings.Contains(string(body), `"ok":true`) {
		t.Fatalf("Do() body = %s, want ok", body)
	}
}

func TestDo_exhaustsRetries(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()

	c := New(Config{MaxAttempts: 3, BaseBackoff: time.Millisecond, MaxBackoff: time.Millisecond})
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
	status, _, err := c.Do(context.Background(), req, nil)
	if !errors.Is(err, ErrTooManyRequests) {
		t.Fatalf("Do() error = %v, want ErrTooManyRequests", err)
	}
	if status != http.StatusTooManyRequests {
		t.Fatalf("Do() status = %d, want 429", status)
	}
}

func TestDo_stopsOnNonRetryableStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	c := New(Config{MaxAttempts: 5, BaseBackoff: time.Millisecond, MaxBackoff: time.Millisecond})
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
	status, _, err := c.Do(context.Background(), req, nil)
	if err == nil {
		t.Fatal("Do() error = nil, want non-nil")
	}
	if status != http.StatusNotFound {
		t.Fatalf("Do() status = %d, want 404", status)
	}
}

func TestDo_honorsRetryAfter(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	start := time.Now()
	c := New(Config{MaxAttempts: 2, BaseBackoff: time.Hour, MaxBackoff: time.Hour})
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
	status, _, err := c.Do(context.Background(), req, nil)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("Do() status = %d, want 200", status)
	}
	// even though base backoff is an hour, Retry-After of 1s must win.
	if elapsed := time.Since(start); elapsed < time.Second || elapsed > 5*time.Second {
		t.Fatalf("Do() elapsed = %s, want ~1s", elapsed)
	}
}

func TestDo_rotatesSessionPerAttempt(t *testing.T) {
	var calls atomic.Int32
	var seen []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.UserAgent()+"|"+r.Header.Get("Cookie"))
		if calls.Add(1) <= 2 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := New(Config{MaxAttempts: 3, BaseBackoff: time.Millisecond, MaxBackoff: time.Millisecond})
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
	if _, _, err := c.Do(context.Background(), req, nil); err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if len(seen) != 3 {
		t.Fatalf("seen = %d, want 3 attempts", len(seen))
	}
	if seen[0] == seen[1] {
		t.Fatalf("session did not rotate: %q", seen[0])
	}
}

func TestDo_contextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cancel()
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()

	c := New(Config{MaxAttempts: 5, BaseBackoff: time.Hour, MaxBackoff: time.Hour})
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL, nil)
	if _, _, err := c.Do(ctx, req, nil); err == nil {
		t.Fatal("Do() error = nil, want context error")
	}
}

func TestDo_injectsSteamHeaders(t *testing.T) {
	var gotUserAgent string
	var gotCookie string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUserAgent = r.UserAgent()
		gotCookie = r.Header.Get("Cookie")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := New(Config{MaxAttempts: 1})
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL+"?steam_id=76561198088587178", nil)
	if _, _, err := c.Do(context.Background(), req, nil); err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if gotUserAgent == "" || !strings.Contains(gotUserAgent, "Mozilla") {
		t.Fatalf("user agent = %q, want browser-like", gotUserAgent)
	}
	if !strings.Contains(gotCookie, "sessionid=") {
		t.Fatalf("cookie = %q, want sessionid", gotCookie)
	}
}
