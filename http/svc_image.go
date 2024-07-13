package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	dgx "github.com/kudarap/dotagiftx"
)

func handleImageUpload(svc dgx.ImageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get uploaded file.
		form, _, err := r.FormFile("file")
		if err != nil {
			respondError(w, fmt.Errorf("could not find 'file' on form-data: %s", err))
			return
		}
		defer form.Close()

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

func handleImage(svc dgx.ImageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		path, err := svc.Image(id)
		if err != nil {
			respondError(w, err)
			return
		}

		cc := fmt.Sprintf("max-age=%d, public", imageCacheMaxAge)
		w.Header().Add("Cache-Control", cc)
		http.ServeFile(w, r, path)
	}
}

func handleImageThumbnail(svc dgx.ImageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		width, _ := strconv.Atoi(chi.URLParam(r, "w"))
		height, _ := strconv.Atoi(chi.URLParam(r, "h"))

		path, err := svc.Thumbnail(id, uint(width), uint(height))
		if err != nil {
			respondError(w, err)
			return
		}

		cc := fmt.Sprintf("max-age=%d, public", imageCacheItemMaxAge)
		w.Header().Add("Cache-Control", cc)
		http.ServeFile(w, r, path)
	}
}
