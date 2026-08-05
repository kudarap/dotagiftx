package steamcrawl

import (
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"
)

// computeDelay returns how long to wait before the next retry attempt.
// A server provided retryAfter (seconds) takes precedence over the
// exponential backoff. Otherwise it returns base * 2^attempt capped at
// max with full jitter added.
func computeDelay(attempt int, retryAfter time.Duration, base, maxBackoff time.Duration) time.Duration {
	if retryAfter > 0 {
		return retryAfter
	}

	backoff := base * time.Duration(1<<max(attempt-1, 0))
	if backoff > maxBackoff {
		backoff = maxBackoff
	}

	// full jitter between 0 and backoff
	return time.Duration(rand.Int64N(int64(backoff)))
}

// parseRetryAfter parses the value of a Retry-After header which can be
// either a number of seconds or a HTTP-date. It returns 0 when unparseable.
func parseRetryAfter(v string) time.Duration {
	if v == "" {
		return 0
	}

	if sec, err := strconv.Atoi(v); err == nil {
		return time.Duration(sec) * time.Second
	}

	if t, err := http.ParseTime(v); err == nil {
		d := time.Until(t)
		if d < 0 {
			return 0
		}
		return d
	}

	return 0
}
