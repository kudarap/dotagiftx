package http

import (
	"log"
	"net/http"
	"reflect"
	"time"

	dgx "github.com/kudarap/dotagiftx"
)

func hydrateStatsMarketSummaryX(cacheKey string, svc dgx.StatsService, cache dgx.Cache) {
	filter := &dgx.Market{Type: dgx.MarketTypeAsk}
	asks, err := svc.CountMarketStatus(dgx.FindOpts{Filter: filter})
	if err != nil {
		return
	}

	filter.Type = dgx.MarketTypeBid
	bids, err := svc.CountMarketStatus(dgx.FindOpts{Filter: filter})
	if err != nil {
		return
	}

	res := struct {
		*dgx.MarketStatusCount
		Bids *dgx.MarketStatusCount `json:"bids"`
	}{asks, bids}

	if err = cache.Set(cacheKey, res, 0); err != nil {
		log.Println("Error hydrateStatsMarketSummaryX", err)
	}
}

func handleStatsMarketSummary(svc dgx.StatsService, cache dgx.Cache) http.HandlerFunc {
	const cacheKeyX = "stats_market_summary_exp"

	go func() {
		t := time.NewTicker(time.Hour / 2)
		for {
			<-t.C
			hydrateStatsMarketSummaryX(cacheKeyX, svc, cache)
		}
	}()

	if hit, _ := cache.Get(cacheKeyX); hit == "" {
		go hydrateStatsMarketSummaryX(cacheKeyX, svc, cache)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		//cacheKey, noCache := core.CacheKeyFromRequestWithPrefix(r, marketCacheKeyPrefix)
		cacheKey, noCache := dgx.CacheKeyFromRequest(r)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		f := &dgx.Market{}
		if err := findOptsFilter(r.URL, f); err != nil {
			respondError(w, err)
			return
		}
		// Use hydration when getting all market status
		if reflect.DeepEqual(f, &dgx.Market{}) {
			hit, _ := cache.Get(cacheKeyX)
			if hit == "" {
				respondOK(w, struct {
					*dgx.MarketStatusCount
					Bids *dgx.MarketStatusCount `json:"bids"`
				}{})
				return
			}
			respondOK(w, hit)
			return
		}

		var err error
		var asks *dgx.MarketStatusCount
		var bids *dgx.MarketStatusCount

		// check for user mode
		if f.UserID != "" {
			stats, errStat := svc.CountUserMarketStatus(f.UserID)
			if errStat != nil {
				respondError(w, errStat)
				return
			}
			asks = stats
			bids = &dgx.MarketStatusCount{
				BidLive:      stats.BidLive,
				BidCompleted: stats.BidCompleted,
			}
		} else {
			f.Type = dgx.MarketTypeAsk
			asks, err = svc.CountMarketStatus(dgx.FindOpts{Filter: f})
			if err != nil {
				respondError(w, err)
				return
			}
			f.Type = dgx.MarketTypeBid
			bids, err = svc.CountMarketStatus(dgx.FindOpts{Filter: f})
			if err != nil {
				respondError(w, err)
				return
			}
		}

		res := struct {
			*dgx.MarketStatusCount
			Bids *dgx.MarketStatusCount `json:"bids"`
		}{asks, bids}

		go cache.Set(cacheKey, res, time.Hour)
		respondOK(w, res)
	}
}

func handleGraphMarketSales(svc dgx.StatsService, cache dgx.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		//cacheKey, noCache := core.CacheKeyFromRequestWithPrefix(r, marketCacheKeyPrefix)
		cacheKey, noCache := dgx.CacheKeyFromRequest(r)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		f := &dgx.Market{}
		if err := findOptsFilter(r.URL, f); err != nil {
			respondError(w, err)
			return
		}

		res, err := svc.GraphMarketSales(dgx.FindOpts{Filter: f})
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

func handleStatsTopOrigins(itemSvc dgx.ItemService, cache dgx.Cache) http.HandlerFunc {
	return topStatsBaseHandler(itemSvc.TopOrigins, cache)
}

func handleStatsTopHeroes(itemSvc dgx.ItemService, cache dgx.Cache) http.HandlerFunc {
	return topStatsBaseHandler(itemSvc.TopHeroes, cache)
}

func handleStatsTopKeywords(statsSvc dgx.StatsService, cache dgx.Cache) http.HandlerFunc {
	const expiration = time.Hour * 12
	return func(w http.ResponseWriter, r *http.Request) {
		cacheKey, noCache := dgx.CacheKeyFromRequest(r)
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

func topStatsBaseHandler(fn func() ([]string, error), cache dgx.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		cacheKey, noCache := dgx.CacheKeyFromRequest(r)
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
