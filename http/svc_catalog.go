package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/kudarap/dotagiftx/core"
	"github.com/sirupsen/logrus"
)

const (
	queryFlagRecentItems    = "recent"
	queryFlagPopularItems   = "popular"
	queryFlagRecentBidItems = "recent-bid"
)

func handleMarketCatalogList(
	svc core.MarketService,
	trackSvc core.TrackService,
	cache core.Cache,
	logger *logrus.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var noCache bool
		query := r.URL.Query()

		// Special query flags with findOpts override for popular and recent items.
		if hasQueryField(r.URL, "sort") {
			switch query.Get("sort") {
			case queryFlagRecentItems:
				query.Set("sort", "recent_ask:desc")
				noCache = true
				break
			case queryFlagPopularItems:
				query.Set("sort", "view_count:desc")
				break
			case queryFlagRecentBidItems:
				query.Set("sort", "recent_bid:desc")
				break
			}

			r.URL.RawQuery = query.Encode()
		}

		opts, err := findOptsFromURL(r.URL, &core.Catalog{})
		if err != nil {
			respondError(w, err)
			return
		}

		go func() {
			if err := trackSvc.CreateSearchKeyword(r, opts.Keyword); err != nil {
				logger.Errorf("search keyword tracking error: %s", err)
			}
		}()

		// Check for cache hit and render them.
		cacheKey, noCache := core.CacheKeyFromRequestWithPrefix(r, marketCacheKeyPrefix)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		list, md, err := svc.Catalog(opts)
		if err != nil {
			respondError(w, err)
			return
		}
		if list == nil {
			list = []core.Catalog{}
		}

		// Save result to cache.
		data := newDataWithMeta(list, md)
		go func() {
			if err := cache.Set(cacheKey, data, marketCacheExpr); err != nil {
				logger.Errorf("could not save cache on catalog list: %s", err)
			}
		}()

		respondOK(w, data)
	}
}

func handleMarketCatalogDetail(svc core.MarketService, cache core.Cache, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		cacheKey, noCache := core.CacheKeyFromRequestWithPrefix(r, marketCacheKeyPrefix)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		c, err := svc.CatalogDetails(chi.URLParam(r, "slug"))
		if err != nil {
			respondError(w, err)
			return
		}

		go func() {
			if err := cache.Set(cacheKey, c, marketCacheExpr); err != nil {
				logger.Errorf("could not save cache on catalog details: %s", err)
			}
		}()

		respondOK(w, c)
	}
}

const catalogTrendCacheExpr = time.Hour * 2

// TODO! this is hotfixed for slow query on trending catalog.
const catalogTrendRehydrationDur = catalogTrendCacheExpr / 2

var catalogTrendLastUpdated = time.Now().Add(catalogTrendRehydrationDur)

func rehydrateCatalogTrend(cacheKey string, svc core.MarketService, cache core.Cache, logger *logrus.Logger) {
	if time.Now().Before(catalogTrendLastUpdated) {
		return
	}
	catalogTrendLastUpdated = time.Now().Add(catalogTrendRehydrationDur)

	logger.Infoln("REHYDRATING...")
	l, _, _ := svc.TrendingCatalog(core.FindOpts{})
	d := newDataWithMeta(l, &core.FindMetadata{len(l), 10})
	if err := cache.Set(cacheKey, d, catalogTrendCacheExpr); err != nil {
		logger.Errorf("could not save cache on catalog trend list: %s", err)
	}
	logger.Infoln("REHYDRATED", d.ResultCount)
}

func handleMarketCatalogTrendList(svc core.MarketService, cache core.Cache, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var noCache bool
		opts, err := findOptsFromURL(r.URL, &core.Catalog{})
		if err != nil {
			respondError(w, err)
			return
		}

		// Check for cache hit and render them.
		cacheKey, noCache := core.CacheKeyFromRequest(r)
		if !noCache {
			// HOTFIXED! rehydrate before cache expiration.
			go rehydrateCatalogTrend(cacheKey, svc, cache, logger)

			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		list, md, err := svc.TrendingCatalog(opts)
		if err != nil {
			respondError(w, err)
			return
		}
		if list == nil {
			list = []core.Catalog{}
		}

		// Save result to cache.
		data := newDataWithMeta(list, md)
		go func() {
			if err := cache.Set(cacheKey, data, catalogTrendCacheExpr); err != nil {
				logger.Errorf("could not save cache on catalog trend list: %s", err)
			}
			catalogTrendLastUpdated = time.Now().Add(catalogTrendRehydrationDur)
		}()

		respondOK(w, data)
	}
}
