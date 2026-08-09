package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/kudarap/dotagiftx"
)

const (
	itemImportFileType = "text/yaml"

	itemCacheKeyPrefix = "svc_item"
	itemCacheExpr      = time.Hour * 24 * 365 // Full year expiration since item update only happens during BP.
)

// itemService provides access to item service methods used by http handlers.
type itemService interface {
	// Items returns a list of items.
	Items(ctx context.Context, opts dotagiftx.FindOpts) ([]dotagiftx.Item, *dotagiftx.FindMetadata, error)

	// Item returns item details by id.
	Item(ctx context.Context, id string) (*dotagiftx.Item, error)

	// Create saves new item details.
	Create(context.Context, *dotagiftx.Item) error

	// Import creates new item from yaml format.
	Import(ctx context.Context, f io.Reader) (dotagiftx.ItemImportResult, error)

	// TopOrigins returns a list of top origin/treasure base on view count.
	TopOrigins(ctx context.Context) ([]string, error)

	// TopHeroes returns a list of top heroes base on view count.
	TopHeroes(ctx context.Context) ([]string, error)
}

func handleItemList(
	svc itemService,
	trackSvc trackService,
	cache cacheManager,
	logger *slog.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		cacheKey, noCache := cacheKeyFromRequestWithPrefix(r, itemCacheKeyPrefix)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		opts, err := findOptsFromURL(r.URL, &dotagiftx.Item{})
		if err != nil {
			respondError(w, err)
			return
		}

		go func() {
			if err = trackSvc.CreateSearchKeyword(context.WithoutCancel(r.Context()), r, opts.Keyword); err != nil {
				logger.ErrorContext(r.Context(), "search keyword tracking error", "error", err)
			}
		}()

		list, md, err := svc.Items(r.Context(), opts)
		if err != nil {
			respondError(w, err)
			return
		}
		if list == nil {
			list = []dotagiftx.Item{}
		}

		o := newDataWithMeta(list, md)
		go func() {
			if err = cache.Set(cacheKey, o, itemCacheExpr); err != nil {
				logger.ErrorContext(r.Context(), "could not save cache on catalog details", "error", err)
			}
		}()
		respondOK(w, o)
	}
}

func handleItemDetail(svc itemService, cache cacheManager, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		cacheKey, noCache := cacheKeyFromRequestWithPrefix(r, itemCacheKeyPrefix)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		i, err := svc.Item(r.Context(), chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, err)
			return
		}

		go func() {
			if err := cache.Set(cacheKey, i, itemCacheExpr); err != nil {
				logger.ErrorContext(r.Context(), "could not save cache on catalog details", "error", err)
			}
		}()
		respondOK(w, i)
	}
}

func handleItemCreate(svc itemService, cache cacheManager, divineKey string, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := isValidDivineKey(r, divineKey); err != nil {
			respondError(w, err)
			return
		}

		i := new(dotagiftx.Item)
		if err := parseForm(r, i); err != nil {
			respondError(w, err)
			return
		}

		if err := svc.Create(r.Context(), i); err != nil {
			respondError(w, err)
			return
		}

		go func() {
			if err := cache.BulkDel(itemCacheKeyPrefix); err != nil {
				logger.ErrorContext(r.Context(), "could not invalidate item cache", "error", err)
			}
		}()

		respondOK(w, i)
	}
}

func handleItemImport(svc itemService, cache cacheManager, divineKey string, logger *slog.Logger) http.HandlerFunc {
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
		defer func() {
			if err := f.Close(); err != nil {
				logger.ErrorContext(r.Context(), "closing import file", "error", err)
			}
		}()

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

		go func() {
			if err := cache.BulkDel(itemCacheKeyPrefix); err != nil {
				logger.ErrorContext(r.Context(), "could not invalidate item cache", "error", err)
			}
		}()

		respondOK(w, res)
	}
}
