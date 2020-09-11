package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/kudarap/dotagiftx/core"
)

func handleImageUpload(svc core.ImageService) http.HandlerFunc {
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

func handleImage(svc core.ImageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		path, err := svc.Image(id)
		if err != nil {
			respondError(w, err)
			return
		}

		http.ServeFile(w, r, path)
	}
}

func handleImageThumbnail(svc core.ImageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		width, _ := strconv.Atoi(chi.URLParam(r, "w"))
		height, _ := strconv.Atoi(chi.URLParam(r, "h"))

		path, err := svc.Thumbnail(id, uint(width), uint(height))
		if err != nil {
			respondError(w, err)
			return
		}

		http.ServeFile(w, r, path)
	}
}
