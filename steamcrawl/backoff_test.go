package steamcrawl

import (
	"net/http"
	"testing"
	"time"
)

func Test_computeDelay_retryAfterWins(t *testing.T) {
	got := computeDelay(1, 10*time.Second, time.Hour, time.Hour)
	if got != 10*time.Second {
		t.Fatalf("computeDelay() = %s, want 10s", got)
	}
}

func Test_computeDelay_growsAndCaps(t *testing.T) {
	base := time.Second
	for attempt := 1; attempt <= 6; attempt++ {
		got := computeDelay(attempt, 0, base, 5*time.Second)
		if got < 0 || got > 5*time.Second {
			t.Fatalf("attempt %d: computeDelay() = %s out of bounds", attempt, got)
		}
	}
	// attempt 3 => 1s * 2^2 = 4s (full jitter 0..4s)
	if got := computeDelay(3, 0, base, time.Hour); got >= 5*time.Second {
		t.Fatalf("attempt 3: computeDelay() = %s, want < 4s", got)
	}
}

func Test_computeDelay_zeroAttempt(t *testing.T) {
	got := computeDelay(0, 0, time.Second, time.Hour)
	if got >= time.Second {
		t.Fatalf("computeDelay(0) = %s, want < 1s", got)
	}
}
func Test_parseRetryAfter(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want time.Duration
	}{
		{"seconds", "5", 5 * time.Second},
		{"http date past", time.Now().Add(-time.Hour).UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT"), 0},
		{"garbage", "nope", 0},
		{"empty", "", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseRetryAfter(tt.in)
			if got != tt.want {
				t.Fatalf("parseRetryAfter() = %s, want %s", got, tt.want)
			}
		})
	}

	t.Run("http date future", func(t *testing.T) {
		in := time.Now().Add(3 * time.Second).UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
		got := parseRetryAfter(in)
		if got < 2*time.Second || got > 4*time.Second {
			t.Fatalf("parseRetryAfter() = %s, want ~3s", got)
		}
	})
}

func Test_retryAfterHeader(t *testing.T) {
	h := http.Header{}
	h.Set("Retry-After", "2")
	if got := retryAfterHeader(h); got != 2*time.Second {
		t.Fatalf("Retry-After = %s, want 2s", got)
	}

	h2 := http.Header{}
	h2.Set("RateLimit-Reset", "3")
	if got := retryAfterHeader(h2); got != 3*time.Second {
		t.Fatalf("RateLimit-Reset = %s, want 3s", got)
	}
}
