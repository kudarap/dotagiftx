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
		s, err := svc.Market(r.Context(), chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, s)
	}
}

func handleMarketCreate(svc core.MarketService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := new(core.Market)
		if err := parseForm(r, s); err != nil {
			respondError(w, err)
			return
		}

		if err := svc.Create(r.Context(), s); err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, s)
	}
}

func handleMarketUpdate(svc core.MarketService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := new(core.Market)
		if err := parseForm(r, s); err != nil {
			respondError(w, err)
			return
		}
		s.ID = chi.URLParam(r, "id")

		if err := svc.Update(r.Context(), s); err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, s)
	}
}

const cacheExpr = time.Minute * 2

func handleMarketIndexList(svc core.MarketService, trackSvc core.TrackService, cache core.Cache, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		list, md, err := svc.Index(opts)
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
			if err := cache.Set(cacheKey, data, cacheExpr); err != nil {
				logger.Errorf("could save cache on market index list: %s", err)
			}
		}()

		respondOK(w, data)
	}
}
