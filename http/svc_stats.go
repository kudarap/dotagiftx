package http

import (
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/kudarap/dotagiftx"
	"github.com/sirupsen/logrus"
)

func hydrateStatsMarketSummaryCache(cacheKey string, svc dotagiftx.StatsService, cache cacheManager) {
	filter := &dotagiftx.Market{Type: dotagiftx.MarketTypeAsk}
	asks, err := svc.CountMarketStatus(dotagiftx.FindOpts{Filter: filter})
	if err != nil {
		return
	}

	filter.Type = dotagiftx.MarketTypeBid
	bids, err := svc.CountMarketStatus(dotagiftx.FindOpts{Filter: filter})
	if err != nil {
		return
	}

	res := struct {
		*dotagiftx.MarketStatusCount
		Bids *dotagiftx.MarketStatusCount `json:"bids"`
	}{asks, bids}

	if err = cache.Set(cacheKey, res, 0); err != nil {
		log.Println("Error hydrateStatsMarketSummaryCache", err)
	}
}

func handleStatsMarketSummary(svc dotagiftx.StatsService, cache cacheManager) http.HandlerFunc {
	const cacheKeyX = "stats_market_summary_exp"

	go func() {
		t := time.NewTicker(time.Hour / 2)
		for {
			<-t.C
			hydrateStatsMarketSummaryCache(cacheKeyX, svc, cache)
		}
	}()

	if hit, _ := cache.Get(cacheKeyX); hit == "" {
		go hydrateStatsMarketSummaryCache(cacheKeyX, svc, cache)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		cacheKey, noCache := cacheKeyFromRequest(r)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		filter := &dotagiftx.Market{}
		if err := findOptsFilter(r.URL, filter); err != nil {
			respondError(w, err)
			return
		}
		// Use hydration when getting all market status
		if reflect.DeepEqual(filter, &dotagiftx.Market{}) {
			hit, _ := cache.Get(cacheKeyX)
			if hit == "" {
				respondOK(w, struct {
					*dotagiftx.MarketStatusCount
					Bids *dotagiftx.MarketStatusCount `json:"bids"`
				}{})
				return
			}
			respondOK(w, hit)
			return
		}

		var err error
		var asks *dotagiftx.MarketStatusCount
		var bids *dotagiftx.MarketStatusCount

		// check for user mode
		switch {
		case filter.UserID != "":
			stats, errStat := svc.CountUserMarketStatus(filter.UserID)
			if errStat != nil {
				respondError(w, errStat)
				return
			}
			asks = stats
			bids = &dotagiftx.MarketStatusCount{
				BidLive:      stats.BidLive,
				BidCompleted: stats.BidCompleted,
			}

		case filter.PartnerSteamID != "":
			stats, errStat := svc.CountUserMarketStatusBySteamID(filter.PartnerSteamID)
			if errStat != nil {
				respondError(w, errStat)
				return
			}
			asks = stats
			bids = &dotagiftx.MarketStatusCount{
				BidLive:      stats.BidLive,
				BidCompleted: stats.BidCompleted,
			}

		default:
			filter.Type = dotagiftx.MarketTypeAsk
			asks, err = svc.CountMarketStatus(dotagiftx.FindOpts{Filter: filter})
			if err != nil {
				respondError(w, err)
				return
			}
			filter.Type = dotagiftx.MarketTypeBid
			bids, err = svc.CountMarketStatus(dotagiftx.FindOpts{Filter: filter})
			if err != nil {
				respondError(w, err)
				return
			}
		}

		res := struct {
			*dotagiftx.MarketStatusCount
			Bids *dotagiftx.MarketStatusCount `json:"bids"`
		}{asks, bids}

		go cache.Set(cacheKey, res, time.Hour)
		respondOK(w, res)
	}
}

func handleStatsMarketSummaryV2(svc dotagiftx.StatsService, cache cacheManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		cacheKey, noCache := cacheKeyFromRequest(r)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		res, err := collectMarketStats(svc, r)
		if err != nil {
			respondError(w, err)
			return
		}

		go cache.Set(cacheKey, res, time.Minute*5)
		respondOK(w, res)
	}
}

const (
	overallStatsCacheExpr      = time.Minute * 30
	overallStatsRehydrationDur = overallStatsCacheExpr / 2
)

func hydrateStatsMarketSummaryOverall(cacheKey string, svc dotagiftx.StatsService, cache cacheManager, logger *logrus.Logger) {
	logger.Infoln("REHYDRATING OVERALL STATS: started")
	res, err := collectMarketStats(svc, nil)
	if err != nil {
		logger.Errorf("REHYDRATING OVERALL STATS: could not get overall market stats: %s", err)
		return
	}

	if err = cache.Set(cacheKey, res, overallStatsCacheExpr); err != nil {
		logger.Errorf("REHYDRATING OVERALL STATS: could not save cache on overall market stats: %s", err)
		return
	}
	logger.Infoln("REHYDRATING OVERALL STATS: completed")
}

func handleStatsMarketSummaryOverall(svc dotagiftx.StatsService, cache cacheManager, logger *logrus.Logger) http.HandlerFunc {
	const cacheKey = "stats_market_summary_overall"

	// hydration setup since this is a long-running process
	go func() {
		t := time.NewTicker(overallStatsRehydrationDur)
		for {
			<-t.C
			hydrateStatsMarketSummaryOverall(cacheKey, svc, cache, logger)
		}
	}()
	if hit, _ := cache.Get(cacheKey); hit == "" {
		go hydrateStatsMarketSummaryOverall(cacheKey, svc, cache, logger)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if hit, _ := cache.Get(cacheKey); hit != "" {
			respondOK(w, hit)
			return
		}

		res := marketStats{&dotagiftx.MarketStatusCount{}, &dotagiftx.MarketStatusCount{}}
		respondOK(w, res)
	}
}

func handleGraphMarketSales(svc dotagiftx.StatsService, cache cacheManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		cacheKey, noCache := cacheKeyFromRequest(r)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		f := &dotagiftx.Market{}
		if err := findOptsFilter(r.URL, f); err != nil {
			respondError(w, err)
			return
		}

		res, err := svc.GraphMarketSales(dotagiftx.FindOpts{Filter: f})
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

func handleStatsTopOrigins(itemSvc dotagiftx.ItemService, cache cacheManager) http.HandlerFunc {
	return topStatsBaseHandler(itemSvc.TopOrigins, cache)
}

func handleStatsTopHeroes(itemSvc dotagiftx.ItemService, cache cacheManager) http.HandlerFunc {
	return topStatsBaseHandler(itemSvc.TopHeroes, cache)
}

func handleStatsTopKeywords(statsSvc dotagiftx.StatsService, cache cacheManager) http.HandlerFunc {
	const expiration = time.Hour * 12
	return func(w http.ResponseWriter, r *http.Request) {
		cacheKey, noCache := cacheKeyFromRequest(r)
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

func topStatsBaseHandler(fn func() ([]string, error), cache cacheManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		cacheKey, noCache := cacheKeyFromRequest(r)
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

type marketStats struct {
	*dotagiftx.MarketStatusCount
	Bids *dotagiftx.MarketStatusCount `json:"bids"`
}

// newMarketStats aggregate market sell and buy stats
// TODO: this should move to service layer.
func newMarketStats(asks *dotagiftx.MarketStatusCount, bids *dotagiftx.MarketStatusCount) *marketStats {
	asks.BidLive = bids.BidLive
	asks.BidCompleted = bids.BidCompleted
	return &marketStats{asks, bids}
}

func collectMarketStats(svc dotagiftx.StatsService, r *http.Request) (*marketStats, error) {
	var err error
	opts := [2]dotagiftx.FindOpts{
		{Filter: &dotagiftx.Market{Type: dotagiftx.MarketTypeAsk}},
		{Filter: &dotagiftx.Market{Type: dotagiftx.MarketTypeBid}},
	}
	if r != nil {
		opts[0], err = findOptsFromURL(r.URL, &dotagiftx.Market{Type: dotagiftx.MarketTypeAsk})
		if err != nil {
			return nil, err
		}
		opts[1], err = findOptsFromURL(r.URL, &dotagiftx.Market{Type: dotagiftx.MarketTypeBid})
		if err != nil {
			return nil, err
		}
	}

	// collect market sell stats
	asks, err := svc.CountMarketStatusV2(opts[0])
	if err != nil {
		return nil, err
	}

	// collect market buy stats
	bids, err := svc.CountMarketStatusV2(opts[1])
	if err != nil {
		return nil, err
	}

	return newMarketStats(asks, bids), nil
}
