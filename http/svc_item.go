package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/kudarap/dota2giftables/core"
	"github.com/sirupsen/logrus"
)

func handleItemList(svc core.ItemService, trackSvc core.TrackService, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opts, err := findOptsFromURL(r.URL, &core.Item{})
		if err != nil {
			respondError(w, err)
			return
		}

		go func() {
			if err := trackSvc.CreateSearchKeyword(r, opts.Keyword); err != nil {
				logger.Errorf("search keyword tracking error: %s", err)
			}
		}()

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
		i, err := svc.Item(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, i)
	}
}

func handleItemCreate(svc core.ItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		i := new(core.Item)
		if err := parseForm(r, i); err != nil {
			respondError(w, err)
			return
		}

		if err := svc.Create(r.Context(), i); err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, i)
	}
}
