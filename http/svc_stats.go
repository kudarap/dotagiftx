package http

import (
	"net/http"
	"time"

	"github.com/kudarap/dotagiftx/core"
)

const statsCacheExpr = time.Hour

func handleStatsTopOrigins(svc core.ItemService, cache core.Cache) http.HandlerFunc {
	return topStatsBaseHandler(svc.TopOrigins, cache)
}

func handleStatsTopHeroes(svc core.ItemService, cache core.Cache) http.HandlerFunc {
	return topStatsBaseHandler(svc.TopHeroes, cache)
}

func topStatsBaseHandler(fn func() ([]string, error), cache core.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		cacheKey, noCache := core.CacheKeyFromRequest(r)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		l, err := fn()
		if err != nil {
			respondError(w, err)
			return
		}
		top10 := l[:10]

		_ = cache.Set(cacheKey, top10, statsCacheExpr)
		respondOK(w, top10)
	}
}
