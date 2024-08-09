package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	dgx "github.com/kudarap/dotagiftx"
	"github.com/sirupsen/logrus"
)

const (
	queryFlagRecentItems    = "recent"
	queryFlagPopularItems   = "popular"
	queryFlagRecentBidItems = "recent-bid"
)

func handleMarketCatalogList(
	svc dgx.MarketService,
	trackSvc dgx.TrackService,
	cache dgx.Cache,
	logger *logrus.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var noCache bool

		// Special query flags with findOpts override for popular and recent items.
		query := r.URL.Query()
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
		sortQueryModifier(r)

		opts, err := findOptsFromURL(r.URL, &dgx.Catalog{})
		if err != nil {
			respondError(w, err)
			return
		}
		// EXPERIMENTAL
		opts.IndexKey = "item_id"

		go func() {
			if err := trackSvc.CreateSearchKeyword(r, opts.Keyword); err != nil {
				logger.Errorf("search keyword tracking error: %s", err)
			}
		}()

		// Check for cache hit and render them.
		cacheKey, noCache := dgx.CacheKeyFromRequestWithPrefix(r, marketCacheKeyPrefix)
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
			list = []dgx.Catalog{}
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

func handleMarketCatalogDetail(svc dgx.MarketService, cache dgx.Cache, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		cacheKey, noCache := dgx.CacheKeyFromRequestWithPrefix(r, marketCacheKeyPrefix)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		// Special query flags with findOpts
		sortQueryModifier(r)

		opts, err := findOptsFromURL(r.URL, &dgx.Market{})
		if err != nil {
			respondError(w, err)
			return
		}
		// EXPERIMENTAL
		opts.IndexKey = "item_id"

		c, err := svc.CatalogDetails(chi.URLParam(r, "slug"), opts)
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

func rehydrateCatalogTrend(cacheKey string, svc dgx.MarketService, cache dgx.Cache, logger *logrus.Logger) {
	if time.Now().Before(catalogTrendLastUpdated) {
		return
	}
	catalogTrendLastUpdated = time.Now().Add(catalogTrendRehydrationDur)

	logger.Infoln("REHYDRATING...")
	l, _, _ := svc.TrendingCatalog(dgx.FindOpts{})
	d := newDataWithMeta(l, &dgx.FindMetadata{len(l), 10})
	if err := cache.Set(cacheKey, d, catalogTrendCacheExpr); err != nil {
		logger.Errorf("could not save cache on catalog trend list: %s", err)
	}
	logger.Infoln("REHYDRATED", d.ResultCount)
}

func hydrateCatalogTrendX(cacheKey string, svc dgx.MarketService, cache dgx.Cache, logger *logrus.Logger) {
	logger.Infoln("REHYDRATING EXP...")
	list, _, err := svc.TrendingCatalog(dgx.FindOpts{})
	if err != nil {
		logger.Errorf("could not get catalog trend list: %s", err)
		return
	}

	trend := newDataWithMeta(list, &dgx.FindMetadata{len(list), 10})
	if err = cache.Set(cacheKey, trend, 0); err != nil {
		logger.Errorf("could not save cache on catalog trend list: %s", err)
		return
	}
	logger.Infoln("REHYDRATED EXP", trend.ResultCount)
}

func handleMarketCatalogTrendListX(svc dgx.MarketService, cache dgx.Cache, logger *logrus.Logger) http.HandlerFunc {
	const cacheKeyX = "catalog_trend_exp"

	go func() {
		t := time.NewTicker(catalogTrendRehydrationDur)
		for {
			<-t.C
			hydrateCatalogTrendX(cacheKeyX, svc, cache, logger)
		}
	}()

	if hit, _ := cache.Get(cacheKeyX); hit == "" {
		logger.Infoln("no cached catalog trend")
		go hydrateCatalogTrendX(cacheKeyX, svc, cache, logger)
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

func handleMarketCatalogTrendList(svc dgx.MarketService, cache dgx.Cache, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var noCache bool
		opts, err := findOptsFromURL(r.URL, &dgx.Catalog{})
		if err != nil {
			respondError(w, err)
			return
		}

		// Check for cache hit and render them.
		cacheKey, noCache := dgx.CacheKeyFromRequest(r)
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
			list = []dgx.Catalog{}
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
