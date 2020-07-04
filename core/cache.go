package core

import (
	"net/http"
	"time"

	"github.com/kudarap/dota2giftables/gokit/hash"
)

// Cache provides access to cache database.
type Cache interface {
	Set(key string, val interface{}, expr time.Duration) error
	Get(key string) (val string, err error)
}

// CacheKeyFromRequest returns cache key from http request.
func CacheKeyFromRequest(r *http.Request) string {
	return r.URL.Path + ":" + hash.MD5(r.URL.RawQuery)
}
