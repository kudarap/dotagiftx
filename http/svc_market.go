package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kudarap/dotagiftx"
	"github.com/sirupsen/logrus"
)

const (
	marketCacheKeyPrefix = "svc_market"   // For cache invalidation control.
	marketCacheExpr      = time.Hour * 24 // Full day cache since its using on-demand invalidation and caching.
)

func handleMarketList(
	svc dotagiftx.MarketService,
	trackSvc dotagiftx.TrackService,
	cache cacheManager,
	logger *logrus.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Redact buyer details flag from public requests.
		shouldRedactUser := !isReqAuthorized(r)

		// Check for cache hit and render them.
		cacheKey, noCache := cacheKeyFromRequestWithPrefix(r, marketCacheKeyPrefix)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				if shouldRedactUser {
					respondOK(w, redactBuyersFromCache(hit))
					return
				}

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
			list = []dotagiftx.Market{}
		}

		data := newDataWithMeta(list, md)
		if err = cache.Set(cacheKey, data, marketCacheExpr); err != nil {
			logger.Errorf("could not save cache on market list: %s", err)
		}

		if shouldRedactUser {
			data.Data = redactBuyers(list)
		}

		respondOK(w, data)
	}
}

func sortQueryModifier(r *http.Request) {
	if !hasQueryField(r.URL, "sort") {
		return
	}

	query := r.URL.Query()
	switch query.Get("sort") {
	case "best":
		query.Set("sort", "user_rank_score:desc")
	case "recent":
		query.Set("sort", "updated_at:desc")
	case "lowest":
		query.Set("sort", "price")
	case "highest":
		query.Set("sort", "price:desc")
	}

	r.URL.RawQuery = query.Encode()
}

func handleMarketDetail(svc dotagiftx.MarketService, cache cacheManager, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Redact buyer details flag from public requests.
		shouldRedactUser := !isReqAuthorized(r)

		// Check for cache hit and render them.
		cacheKey, noCache := cacheKeyFromRequestWithPrefix(r, marketCacheKeyPrefix)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				if shouldRedactUser {
					respondOK(w, redactBuyerFromCache(hit))
					return
				}

				respondOK(w, hit)
				return
			}
		}

		m, err := svc.Market(r.Context(), chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, err)
			return
		}

		if err := cache.Set(cacheKey, m, marketCacheExpr); err != nil {
			logger.Errorf("could not save cache on market list: %s", err)
		}

		if shouldRedactUser {
			m = redactBuyer(m)
		}

		respondOK(w, m)
	}
}

func handleMarketCreate(svc dotagiftx.MarketService, cache cacheManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := new(dotagiftx.Market)
		if err := parseForm(r, m); err != nil {
			respondError(w, err)
			return
		}

		if err := svc.Create(r.Context(), m); err != nil {
			respondError(w, err)
			return
		}

		go cache.BulkDel(marketCacheKeyPrefix)

		respondOK(w, m)
	}
}

func handleMarketUpdate(svc dotagiftx.MarketService, cache cacheManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := new(dotagiftx.Market)
		if err := parseForm(r, m); err != nil {
			respondError(w, err)
			return
		}
		m.ID = chi.URLParam(r, "id")

		if err := svc.Update(r.Context(), m); err != nil {
			respondError(w, err)
			return
		}

		go cache.BulkDel(marketCacheKeyPrefix)
		respondOK(w, m)
	}
}

func isReqAuthorized(r *http.Request) bool {
	c, _ := ParseFromHeader(r.Header)
	if c == nil {
		return false
	}

	return c.UserID != ""
}

const redactChar = "█"

func redactBuyers(list []dotagiftx.Market) []dotagiftx.Market {
	rl := make([]dotagiftx.Market, len(list))
	copy(rl, list)
	for _, r := range rl {
		if r.Type != dotagiftx.MarketTypeBid {
			continue
		}

		r.User.ID = ""
		r.User.Name = strings.Repeat(redactChar, len(r.User.Name))
		r.User.SteamID = strings.Repeat(redactChar, 10)
		r.User.URL = strings.Repeat(redactChar, 10)
	}

	return rl
}

func redactBuyersFromCache(hit string) interface{} {
	d := struct {
		Data        []dotagiftx.Market `json:"data"`
		ResultCount int                `json:"result_count"`
		TotalCount  int                `json:"total_count"`
	}{}
	if err := json.UnmarshalFromString(hit, &d); err != nil {
		return nil
	}

	d.Data = redactBuyers(d.Data)
	return d
}

func redactBuyer(m *dotagiftx.Market) *dotagiftx.Market {
	if m == nil {
		return nil
	}

	return &redactBuyers([]dotagiftx.Market{*m})[0]
}

func redactBuyerFromCache(hit string) *dotagiftx.Market {
	d := &dotagiftx.Market{}
	if err := json.UnmarshalFromString(hit, &d); err != nil {
		return nil
	}

	return redactBuyer(d)
}
