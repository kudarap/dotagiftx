package core

import (
	"net/http"
	"strings"
	"time"

	"github.com/kudarap/dota2giftables/gokit/hash"
)

// Cache provides access to cache database.
type Cache interface {
	Set(key string, val interface{}, expr time.Duration) error
	Get(key string) (val string, err error)
}

const noCacheReqFlag = "nocache"

// CacheKeyFromRequest returns cache key from http request.
// nocache from request query will return empty string and can be use to skipping cache.
func CacheKeyFromRequest(r *http.Request) string {
	// Skip caching when nocache flag exists.
	if strings.Contains(r.URL.RawQuery, noCacheReqFlag) {
		return ""
	}

	return r.URL.Path + ":" + hash.MD5(r.URL.RawQuery)
}
