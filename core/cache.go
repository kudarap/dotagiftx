package core

import (
	"net/http"
	"strconv"
	"time"

	"github.com/kudarap/dota2giftables/gokit/hash"
)

// Cache provides access to cache database.
type Cache interface {
	Set(key string, val interface{}, expr time.Duration) error
	Get(key string) (val string, err error)
}

const cacheSkipKey = "nocache"

// CacheKeyFromRequest returns cache key from http request.
// nocache from request query will return empty string and can be use to skipping cache.
func CacheKeyFromRequest(r *http.Request) (key string, noCache bool) {
	// Skip caching when nocache flag exists.
	noCache, _ = strconv.ParseBool(r.URL.Query().Get(cacheSkipKey))
	return r.URL.Path + ":" + hash.MD5(r.URL.RawQuery), noCache
}
