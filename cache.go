package dgx

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kudarap/dotagiftx/gokit/hash"
)

// Cache provides access to cache database.
type Cache interface {
	Set(key string, val interface{}, expr time.Duration) error
	Get(key string) (val string, err error)
	Del(key string) error
	BulkDel(keyPrefix string) error
}

const cacheSkipKey = "nocache"

// CacheKeyFromRequest returns cache key from http request.
// nocache from request query will return empty string and can be use to skipping cache.
func CacheKeyFromRequest(r *http.Request) (key string, noCache bool) {
	// Skip caching when nocache flag exists.
	_, noCache = r.URL.Query()[cacheSkipKey]
	// Set owner user id for scoped requests.
	var userID string
	au := AuthFromContext(r.Context())
	if au != nil {
		userID = au.UserID
	}

	// Compose cache key and omit nocache param, this will enable force reloads.
	q := r.URL.Query()
	q.Del(cacheSkipKey)
	key = fmt.Sprintf("%s%s:%s", userID, r.URL.Path, hash.MD5(q.Encode()))
	return
}

func CacheKeyFromRequestWithPrefix(r *http.Request, prefix string) (key string, noCache bool) {
	key, noCache = CacheKeyFromRequest(r)
	key = prefix + ":" + key
	return
}
