package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/kudarap/dota2giftables/core"
)

func handleSellList(svc core.SellService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opts, err := findOptsFromURL(r.URL, &core.Item{})
		if err != nil {
			respondError(w, err)
			return
		}

		list, md, err := svc.Sells(opts)
		if err != nil {
			respondError(w, err)
			return
		}
		if list == nil {
			list = []core.Sell{}
		}

		respondOK(w, newDataWithMeta(list, md))
	}
}

func handleSellDetail(svc core.SellService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, err := svc.Sell(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, s)
	}
}

func handleSellCreate(svc core.SellService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := new(core.Sell)
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

func handleSellUpdate(svc core.SellService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := new(core.Sell)
		if err := parseForm(r, p); err != nil {
			respondError(w, err)
			return
		}
		p.ID = chi.URLParam(r, "id")

		if err := svc.Update(r.Context(), p); err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, p)
	}
}
