package http

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/kudarap/dotagiftx"
)

// cache provides access to cache database.
type cache interface {
	Set(key string, val interface{}, expr time.Duration) error
	Get(key string) (val string, err error)
	BulkDel(keyPrefix string) error
}

const cacheSkipKey = "nocache"

// cacheKeyFromRequest returns cache key from http request.
// nocache from a request query will return empty string and can be used to skipping cache.
func cacheKeyFromRequest(r *http.Request) (key string, noCache bool) {
	// Skip caching when a nocache flag exists.
	_, noCache = r.URL.Query()[cacheSkipKey]
	// Set owner user id for scoped requests.
	var userID string
	au := dotagiftx.AuthFromContext(r.Context())
	if au != nil {
		userID = au.UserID
	}

	// Compose cache key and omit nocache param, this will enable force reloads.
	q := r.URL.Query()
	q.Del(cacheSkipKey)
	h := md5.New()
	h.Write([]byte(q.Encode()))
	hash := hex.EncodeToString(h.Sum(nil))
	key = fmt.Sprintf("%s%s:%s", userID, r.URL.Path, hash)
	return
}

func cacheKeyFromRequestWithPrefix(r *http.Request, prefix string) (key string, noCache bool) {
	key, noCache = cacheKeyFromRequest(r)
	key = prefix + ":" + key
	return
}
