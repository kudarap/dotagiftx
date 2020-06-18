package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/kudarap/dota2giftables/core"
)

func handleItemList(svc core.ItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opts, err := findOptsFromURL(r.URL, &core.Item{})
		if err != nil {
			respondError(w, err)
			return
		}

		list, md, err := svc.Items(opts)
		if err != nil {
			respondError(w, err)
			return
		}
		if list == nil {
			list = []core.Item{}
		}

		respondOK(w, newDataWithMeta(list, md))
	}
}

func handleItemDetail(svc core.ItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		item, err := svc.Item(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, item)
	}
}

func handleItemCreate(svc core.ItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		item := new(core.Item)
		if err := parseForm(r, item); err != nil {
			respondError(w, err)
			return
		}

		if err := svc.Create(r.Context(), item); err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, item)
	}
}
