package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/kudarap/dota2giftables/core"
)

func handleMarketList(svc core.MarketService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opts, err := findOptsFromURL(r.URL, &core.Market{})
		if err != nil {
			respondError(w, err)
			return
		}

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

func handleMarketIndexList(svc core.MarketService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opts, err := findOptsFromURL(r.URL, &core.Market{})
		if err != nil {
			respondError(w, err)
			return
		}

		list, md, err := svc.Index(opts)
		if err != nil {
			respondError(w, err)
			return
		}
		if list == nil {
			list = []core.MarketIndex{}
		}

		respondOK(w, newDataWithMeta(list, md))
	}
}
