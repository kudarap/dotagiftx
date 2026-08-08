package http

import (
	"context"
	"net/http"
	"time"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/kudarap/dotagiftx"
)

const (
	queryFlagRecentItems    = "recent"
	queryFlagPopularItems   = "popular"
	queryFlagRecentBidItems = "recent-bid"
)

func handleMarketCatalogList(
	svc dotagiftx.MarketService,
	trackSvc dotagiftx.TrackService,
	cache cacheManager,
	logger *slog.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var noCache bool

		// Special query flags with findOpts override for popular and recent items.
		query := r.URL.Query()
		if hasQueryField(r.URL, "sort") {
			switch query.Get("sort") {
			case queryFlagRecentItems:
				query.Set("sort", "recent_ask:desc")
			case queryFlagPopularItems:
				query.Set("sort", "view_count:desc")
			case queryFlagRecentBidItems:
				query.Set("sort", "recent_bid:desc")
			}

			r.URL.RawQuery = query.Encode()
		}
		sortQueryModifier(r)

		opts, err := findOptsFromURL(r.URL, &dotagiftx.Catalog{})
		if err != nil {
			respondError(w, err)
			return
		}
		// EXPERIMENTAL
		opts.IndexKey = "item_id"

		go func() {
			if err := trackSvc.CreateSearchKeyword(context.Background(), r, opts.Keyword); err != nil {
				logger.Error("search keyword tracking error", "error", err)
			}
		}()

		// Check for cache hit and render them.
		cacheKey, noCache := cacheKeyFromRequestWithPrefix(r, marketCacheKeyPrefix)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		list, md, err := svc.Catalog(r.Context(), opts)
		if err != nil {
			respondError(w, err)
			return
		}
		if list == nil {
			list = []dotagiftx.Catalog{}
		}

		// Save result to cache.
		data := newDataWithMeta(list, md)
		go func() {
			if err := cache.Set(cacheKey, data, marketCacheExpr); err != nil {
				logger.Error("could not save cache on catalog list", "error", err)
			}
		}()

		respondOK(w, data)
	}
}

func handleMarketCatalogDetail(svc dotagiftx.MarketService, cache cacheManager, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		cacheKey, noCache := cacheKeyFromRequestWithPrefix(r, marketCacheKeyPrefix)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		// Special query flags with findOpts
		sortQueryModifier(r)

		opts, err := findOptsFromURL(r.URL, &dotagiftx.Market{})
		if err != nil {
			respondError(w, err)
			return
		}
		// EXPERIMENTAL
		opts.IndexKey = "item_id"

		c, err := svc.CatalogDetails(r.Context(), chi.URLParam(r, "slug"), opts)
		if err != nil {
			respondError(w, err)
			return
		}

		go func() {
			if err := cache.Set(cacheKey, c, marketCacheExpr); err != nil {
				logger.Error("could not save cache on catalog details", "error", err)
			}
		}()

		respondOK(w, c)
	}
}

const catalogTrendCacheExpr = time.Hour * 2

// TODO! this is hotfixed for slow query on trending catalog.
const catalogTrendRehydrationDur = catalogTrendCacheExpr / 2

func hydrateCatalogTrend(cacheKey string, svc dotagiftx.MarketService, cache cacheManager, logger *slog.Logger) {
	logger.Info("REHYDRATING EXP...")
	list, _, err := svc.TrendingCatalog(context.Background(), dotagiftx.FindOpts{})
	if err != nil {
		logger.Error("could not get catalog trend list", "error", err)
		return
	}

	trend := newDataWithMeta(list, &dotagiftx.FindMetadata{ResultCount: len(list), TotalCount: 10})
	if err = cache.Set(cacheKey, trend, 0); err != nil {
		logger.Error("could not save cache on catalog trend list", "error", err)
		return
	}
	logger.Info("REHYDRATED EXP", "result_count", trend.ResultCount)
}

func handleMarketCatalogTrendList(svc dotagiftx.MarketService, cache cacheManager, logger *slog.Logger) http.HandlerFunc {
	const cacheKeyX = "catalog_trend_exp"

	go func() {
		t := time.NewTicker(catalogTrendRehydrationDur)
		for {
			<-t.C
			hydrateCatalogTrend(cacheKeyX, svc, cache, logger)
		}
	}()

	if hit, _ := cache.Get(cacheKeyX); hit == "" {
		logger.Info("no cached catalog trend")
		go hydrateCatalogTrend(cacheKeyX, svc, cache, logger)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		hit, _ := cache.Get(cacheKeyX)
		if hit == "" {
			hit = `{
    "data": null,
    "result_count": 0,
    "total_count": 10
}`
		}
		respondOK(w, hit)
	}
}
