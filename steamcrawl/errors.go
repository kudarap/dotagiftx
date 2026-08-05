package steamcrawl

import (
	"errors"
	"fmt"
)

var (
	// ErrTooManyRequests is returned when a request was rate limited
	// and retries have been exhausted.
	ErrTooManyRequests = errors.New("too many requests")
	// ErrRetriesExhausted is returned when all retries were exhausted
	// without a successful response.
	ErrRetriesExhausted = errors.New("retries exhausted")
)

type statusError struct {
	statusCode int
	text       string
}

func (e *statusError) Error() string {
	return fmt.Sprintf("unexpected status: %d %s", e.statusCode, e.text)
}
