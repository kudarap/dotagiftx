package http

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-chi/chi"
	"github.com/kudarap/dota2giftables/core"
)

func handleMarketList(svc core.MarketService, trackSvc core.TrackService, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opts, err := findOptsFromURL(r.URL, &core.Market{})
		if err != nil {
			respondError(w, err)
			return
		}

		go func() {
			if err := trackSvc.CreateSearchKeyword(r, opts.Keyword); err != nil {
				logger.Errorf("search keyword tracking error: %s", err)
			}
		}()

		list, md, err := svc.Markets(r.Context(), opts)
		if err != nil {
			respondError(w, err)
			return
		}
		if list == nil {
			list = []core.Market{}
		}

		respondOK(w, newDataWithMeta(list, md))
	}
}

func handleMarketDetail(svc core.MarketService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m, err := svc.Market(r.Context(), chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, m)
	}
}

func handleMarketCreate(svc core.MarketService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := new(core.Market)
		if err := parseForm(r, m); err != nil {
			respondError(w, err)
			return
		}

		if err := svc.Create(r.Context(), m); err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, m)
	}
}

func handleMarketUpdate(svc core.MarketService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := new(core.Market)
		if err := parseForm(r, m); err != nil {
			respondError(w, err)
			return
		}
		m.ID = chi.URLParam(r, "id")

		if err := svc.Update(r.Context(), m); err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, m)
	}
}

const (
	catalogCacheExpr      = time.Minute
	queryFlagRecentItems  = "recent-items"
	queryFlagPopularItems = "popular-items"
)

func handleMarketCatalogList(svc core.MarketService, trackSvc core.TrackService, cache core.Cache, logger *logrus.Logger) http.HandlerFunc {
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
		cacheKey, noCache := core.CacheKeyFromRequest(r)

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
			if err := cache.Set(cacheKey, data, catalogCacheExpr); err != nil {
				logger.Errorf("could save cache on market index list: %s", err)
			}
		}()

		respondOK(w, data)
	}
}

func handleMarketCatalogDetail(svc core.MarketService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := svc.CatalogDetails(chi.URLParam(r, "slug"))
		if err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, c)
	}
}
