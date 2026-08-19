package http

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/kudarap/dotagiftx/assets"
)

const pixelImage = "image/pixel.gif"

// trackService provides access to track service methods used by http handlers.
type trackService interface {
	// CreateFromRequest saves new track from http request. Primarily used on client side.
	CreateFromRequest(ctx context.Context, r *http.Request) error

	// CreateSearchKeyword saves new keyword tracking data.
	CreateSearchKeyword(ctx context.Context, r *http.Request, keyword string) error
}

func handleTracker(svc trackService, logger *slog.Logger) http.HandlerFunc {
	image, _ := assets.Content.ReadFile(pixelImage)
	return func(w http.ResponseWriter, r *http.Request) {
		go func(r *http.Request) {
			if err := svc.CreateFromRequest(context.Background(), r); err != nil {
				logger.ErrorContext(r.Context(), "tracker error", "error", err)
			}
		}(r)

		// unset JSON headers
		w.Header().Set("Access-Control-Allow-Headers", "")
		w.Header().Set("Access-Control-Allow-Methods", "")
		w.Header().Set("Access-Control-Allow-Origin", "")

		// no cache
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		// output image
		w.Header().Set("Content-Type", "image/gif")
		_, _ = w.Write(image)
	}
}
