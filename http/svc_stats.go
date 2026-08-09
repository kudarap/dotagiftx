package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/kudarap/dotagiftx/dotagiftx"
)

// statsService provides access to stats service methods used by http handlers.
type statsService interface {
	// CountMarketStatusV2 returns market status count base on given options.
	CountMarketStatusV2(ctx context.Context, opts dotagiftx.FindOpts) (*dotagiftx.MarketStatusCount, error)

	// GraphMarketSales returns market sales graph base on given options.
	GraphMarketSales(ctx context.Context, opts dotagiftx.FindOpts) ([]dotagiftx.MarketSalesGraph, error)

	// TopKeywords returns a list of top search keywords.
	TopKeywords(ctx context.Context) ([]dotagiftx.SearchKeywordScore, error)
}

func handleStatsMarketSummaryV2(svc statsService, cache cacheManager, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		cacheKey, noCache := cacheKeyFromRequest(r)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		res, err := collectMarketStats(r.Context(), svc, r)
		if err != nil {
			respondError(w, err)
			return
		}

		go func() {
			if err := cache.Set(cacheKey, res, time.Minute*5); err != nil {
				logger.ErrorContext(r.Context(), "could not save cache on market summary", "error", err)
			}
		}()
		respondOK(w, res)
	}
}

const (
	overallStatsCacheExpr      = time.Minute * 30
	overallStatsRehydrationDur = overallStatsCacheExpr / 2
)

func hydrateStatsMarketSummaryOverall(cacheKey string, svc statsService, cache cacheManager, logger *slog.Logger) {
	ctx := context.Background()
	logger.InfoContext(ctx, "REHYDRATING OVERALL STATS: started")
	res, err := collectMarketStats(ctx, svc, nil)
	if err != nil {
		logger.ErrorContext(ctx, "REHYDRATING OVERALL STATS: could not get overall market stats", "error", err)
		return
	}

	if err = cache.Set(cacheKey, res, overallStatsCacheExpr); err != nil {
		logger.ErrorContext(ctx, "REHYDRATING OVERALL STATS: could not save cache on overall market stats", "error", err)
		return
	}
	logger.InfoContext(ctx, "REHYDRATING OVERALL STATS: completed")
}

func handleStatsMarketSummaryOverall(svc statsService, cache cacheManager, logger *slog.Logger) http.HandlerFunc {
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

func handleGraphMarketSales(svc statsService, cache cacheManager, logger *slog.Logger) http.HandlerFunc {
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

		res, err := svc.GraphMarketSales(r.Context(), dotagiftx.FindOpts{Filter: f})
		if err != nil {
			respondError(w, err)
			return
		}

		const expiration = time.Hour * 4
		go func() {
			if err := cache.Set(cacheKey, res, expiration); err != nil {
				logger.ErrorContext(r.Context(), "could not save cache on market sales graph", "error", err)
			}
		}()
		respondOK(w, res)
	}
}

const statsCacheExpr = time.Hour

func handleStatsTopOrigins(itemSvc itemService, cache cacheManager, logger *slog.Logger) http.HandlerFunc {
	return topStatsBaseHandler(itemSvc.TopOrigins, cache, logger)
}

func handleStatsTopHeroes(itemSvc itemService, cache cacheManager, logger *slog.Logger) http.HandlerFunc {
	return topStatsBaseHandler(itemSvc.TopHeroes, cache, logger)
}

func handleStatsTopKeywords(statsSvc statsService, cache cacheManager, logger *slog.Logger) http.HandlerFunc {
	const expiration = time.Hour * 12
	return func(w http.ResponseWriter, r *http.Request) {
		cacheKey, noCache := cacheKeyFromRequest(r)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		res, err := statsSvc.TopKeywords(r.Context())
		if err != nil {
			respondError(w, err)
			return
		}

		go func() {
			if err := cache.Set(cacheKey, res, expiration); err != nil {
				logger.ErrorContext(r.Context(), "could not save cache on top keywords", "error", err)
			}
		}()
		respondOK(w, res)
	}
}

func topStatsBaseHandler(fn func(context.Context) ([]string, error), cache cacheManager, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		cacheKey, noCache := cacheKeyFromRequest(r)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		l, err := fn(r.Context())
		if err != nil {
			respondError(w, err)
			return
		}

		if len(l) <= 10 {
			respondOK(w, l)
			return
		}

		top10 := l[:10]
		go func() {
			if err := cache.Set(cacheKey, top10, statsCacheExpr); err != nil {
				logger.ErrorContext(r.Context(), "could not save cache on top stats", "error", err)
			}
		}()
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

func collectMarketStats(ctx context.Context, svc statsService, r *http.Request) (*marketStats, error) {
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
	asks, err := svc.CountMarketStatusV2(ctx, opts[0])
	if err != nil {
		return nil, err
	}

	// collect market buy stats
	bids, err := svc.CountMarketStatusV2(ctx, opts[1])
	if err != nil {
		return nil, err
	}

	return newMarketStats(asks, bids), nil
}
