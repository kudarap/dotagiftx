package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	dgx "github.com/kudarap/dotagiftx"
	"github.com/sirupsen/logrus"
)

const (
	itemImportFileType = "text/yaml"

	itemCacheKeyPrefix = "svc_item"
	itemCacheExpr      = time.Hour * 24 * 365 // Full year expiration since item update only happens during BP.
)

func handleItemList(
	svc dgx.ItemService,
	trackSvc dgx.TrackService,
	cache dgx.Cache,
	logger *logrus.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		cacheKey, noCache := dgx.CacheKeyFromRequestWithPrefix(r, itemCacheKeyPrefix)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		opts, err := findOptsFromURL(r.URL, &dgx.Item{})
		if err != nil {
			respondError(w, err)
			return
		}

		go func() {
			if err = trackSvc.CreateSearchKeyword(r, opts.Keyword); err != nil {
				logger.Errorf("search keyword tracking error: %s", err)
			}
		}()

		list, md, err := svc.Items(opts)
		if err != nil {
			respondError(w, err)
			return
		}
		if list == nil {
			list = []dgx.Item{}
		}

		o := newDataWithMeta(list, md)
		go func() {
			if err = cache.Set(cacheKey, o, itemCacheExpr); err != nil {
				logger.Errorf("could not save cache on catalog details: %s", err)
			}
		}()
		respondOK(w, o)
	}
}

func handleItemDetail(svc dgx.ItemService, cache dgx.Cache, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		cacheKey, noCache := dgx.CacheKeyFromRequestWithPrefix(r, itemCacheKeyPrefix)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		i, err := svc.Item(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, err)
			return
		}

		go func() {
			if err := cache.Set(cacheKey, i, itemCacheExpr); err != nil {
				logger.Errorf("could not save cache on catalog details: %s", err)
			}
		}()
		respondOK(w, i)
	}
}

func handleItemCreate(svc dgx.ItemService, cache dgx.Cache, divineKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := isValidDivineKey(r, divineKey); err != nil {
			respondError(w, err)
			return
		}

		i := new(dgx.Item)
		if err := parseForm(r, i); err != nil {
			respondError(w, err)
			return
		}

		if err := svc.Create(r.Context(), i); err != nil {
			respondError(w, err)
			return
		}

		go cache.BulkDel(itemCacheKeyPrefix)

		respondOK(w, i)
	}
}

func handleItemImport(svc dgx.ItemService, cache dgx.Cache, divineKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := isValidDivineKey(r, divineKey); err != nil {
			respondError(w, err)
			return
		}

		// Get uploaded file.
		f, fh, err := r.FormFile("file")
		if err != nil {
			respondError(w, fmt.Errorf("could not find 'file' on form-data: %s", err))
			return
		}
		defer f.Close()

		// Check and read yaml file.
		ct := fh.Header.Get("content-type")
		if ct != itemImportFileType {
			respondError(w, fmt.Errorf("could not parse content-type: %s", ct))
			return
		}

		res, err := svc.Import(r.Context(), f)
		if err != nil {
			respondError(w, err)
			return
		}

		go cache.BulkDel(itemCacheKeyPrefix)

		respondOK(w, res)
	}
}
