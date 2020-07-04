package http

import (
	"fmt"
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

func handleItemImport(svc core.ItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get uploaded file.
		f, fh, err := r.FormFile("file")
		if err != nil {
			respondError(w, fmt.Errorf("could not find 'file' on form-data: %s", err))
			return
		}
		defer f.Close()

		// Read yaml file.
		ct := fh.Header.Get("content-type")
		if ct != "text/yaml" {
			respondError(w, fmt.Errorf("could not parse content-type: %s", ct))
			return
		}

		res, err := svc.Import(r.Context(), f)
		if err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, res)
	}
}
