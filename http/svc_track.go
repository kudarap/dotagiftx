package http

import (
	"net/http"

	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/assets"
	"github.com/sirupsen/logrus"
)

const pixelImage = "image/pixel.gif"

func handleTracker(svc dgx.TrackService, logger *logrus.Logger) http.HandlerFunc {
	image, _ := assets.Content.ReadFile(pixelImage)
	return func(w http.ResponseWriter, r *http.Request) {
		go func(r *http.Request) {
			if err := svc.CreateFromRequest(r); err != nil {
				logger.Errorf("tracker error: %s", err)
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
		w.Write(image)
	}
}
