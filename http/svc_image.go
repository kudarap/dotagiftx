package http

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// imageService provides access to image service methods used by http handlers.
type imageService interface {
	// Upload saves image details and actual file to local file system.
	Upload(context.Context, io.Reader) (fileID string, err error)

	// Image returns image details by id.
	Image(ctx context.Context, fileID string) (path string, err error)

	// Thumbnail downscales an image preserving its aspect ratio to the maximum dimensions.
	Thumbnail(ctx context.Context, fileID string, width, height uint) (path string, err error)
}

func handleImageUpload(svc imageService, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get uploaded file.
		form, _, err := r.FormFile("file")
		if err != nil {
			respondError(w, fmt.Errorf("could not find 'file' on form-data: %s", err))
			return
		}
		defer func() {
			if err := form.Close(); err != nil {
				logger.ErrorContext(r.Context(), "closing upload file", "error", err)
			}
		}()

		id, err := svc.Upload(r.Context(), form)
		if err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, struct {
			FileID string `json:"file_id"`
		}{id})
	}
}

const (
	dayAge               = 3600 * 24    // 1 day
	imageCacheMaxAge     = dayAge       // 1 day for profile and raw image
	imageCacheItemMaxAge = dayAge * 365 // 1 year for item images
)

func handleImage(svc imageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		path, err := svc.Image(r.Context(), id)
		if err != nil {
			respondError(w, err)
			return
		}

		cc := fmt.Sprintf("max-age=%d, public", imageCacheMaxAge)
		w.Header().Add("Cache-Control", cc)
		http.ServeFile(w, r, path) //nolint:gosec // path is resolved by file manager within save dir
	}
}

func handleImageThumbnail(svc imageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		width, _ := strconv.Atoi(chi.URLParam(r, "w"))
		height, _ := strconv.Atoi(chi.URLParam(r, "h"))

		path, err := svc.Thumbnail(r.Context(), id, uint(width), uint(height))
		if err != nil {
			respondError(w, err)
			return
		}

		cc := fmt.Sprintf("max-age=%d, public", imageCacheItemMaxAge)
		w.Header().Add("Cache-Control", cc)
		http.ServeFile(w, r, path) //nolint:gosec // path is resolved by file manager within save dir
	}
}
