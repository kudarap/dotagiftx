package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/kudarap/dotagiftx/core"
)

const userCacheExpr = time.Minute * 5

func handleProfile(svc core.UserService, cache core.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		cacheKey, noCache := core.CacheKeyFromRequest(r)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		u, err := svc.UserFromContext(r.Context())
		if err != nil {
			respondError(w, err)
			return
		}

		go cache.Set(cacheKey, u, userCacheExpr)

		respondOK(w, u)
	}
}

func handlePublicProfile(svc core.UserService, cache core.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		cacheKey, noCache := core.CacheKeyFromRequest(r)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		id := chi.URLParam(r, "id")
		u, err := svc.User(id)
		if err != nil {
			respondError(w, err)
			return
		}

		go cache.Set(cacheKey, u, userCacheExpr)

		respondOK(w, u)
	}
}
