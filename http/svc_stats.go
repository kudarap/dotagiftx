package http

import (
	"net/http"
	"time"

	"github.com/kudarap/dotagiftx/core"
)

func handleStatsMarketSummary(svc core.StatsService, cache core.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		//cacheKey, noCache := core.CacheKeyFromRequestWithPrefix(r, marketCacheKeyPrefix)
		cacheKey, noCache := core.CacheKeyFromRequest(r)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		f := &core.Market{}
		if err := findOptsFilter(r.URL, f); err != nil {
			respondError(w, err)
			return
		}

		var err error
		var asks *core.MarketStatusCount
		var bids *core.MarketStatusCount

		// check for user mode
		if f.UserID != "" {
			stats, errStat := svc.CountUserMarketStatus(f.UserID)
			if errStat != nil {
				respondError(w, errStat)
				return
			}
			asks = stats
			bids = &core.MarketStatusCount{
				BidLive:      stats.BidLive,
				BidCompleted: stats.BidCompleted,
			}
		} else {
			f.Type = core.MarketTypeAsk
			asks, err = svc.CountMarketStatus(core.FindOpts{Filter: f})
			if err != nil {
				respondError(w, err)
				return
			}
			f.Type = core.MarketTypeBid
			bids, err = svc.CountMarketStatus(core.FindOpts{Filter: f})
			if err != nil {
				respondError(w, err)
				return
			}
		}

		res := struct {
			*core.MarketStatusCount
			Bids *core.MarketStatusCount `json:"bids"`
		}{asks, bids}

		go cache.Set(cacheKey, res, time.Hour)
		respondOK(w, res)
	}
}

func handleGraphMarketSales(svc core.StatsService, cache core.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		//cacheKey, noCache := core.CacheKeyFromRequestWithPrefix(r, marketCacheKeyPrefix)
		cacheKey, noCache := core.CacheKeyFromRequest(r)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		f := &core.Market{}
		if err := findOptsFilter(r.URL, f); err != nil {
			respondError(w, err)
			return
		}

		res, err := svc.GraphMarketSales(core.FindOpts{Filter: f})
		if err != nil {
			respondError(w, err)
			return
		}

		const expiration = time.Hour * 4
		go cache.Set(cacheKey, res, expiration)
		respondOK(w, res)
	}
}

const statsCacheExpr = time.Hour

func handleStatsTopOrigins(itemSvc core.ItemService, cache core.Cache) http.HandlerFunc {
	return topStatsBaseHandler(itemSvc.TopOrigins, cache)
}

func handleStatsTopHeroes(itemSvc core.ItemService, cache core.Cache) http.HandlerFunc {
	return topStatsBaseHandler(itemSvc.TopHeroes, cache)
}

func handleStatsTopKeywords(statsSvc core.StatsService, cache core.Cache) http.HandlerFunc {
	const expiration = time.Hour * 12
	return func(w http.ResponseWriter, r *http.Request) {
		cacheKey, noCache := core.CacheKeyFromRequest(r)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		res, err := statsSvc.TopKeywords()
		if err != nil {
			respondError(w, err)
			return
		}

		go cache.Set(cacheKey, res, expiration)
		respondOK(w, res)
	}
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

		go cache.Set(cacheKey, top10, statsCacheExpr)
		respondOK(w, top10)
	}
}
