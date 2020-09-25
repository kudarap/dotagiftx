package http

import (
	"net/http"

	"github.com/kudarap/dotagiftx/core"
)

func handleStatsTopOrigins(svc core.ItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, err := svc.TopOrigins()
		if err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, l[:10])
	}
}

func handleStatsTopHeroes(svc core.ItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, err := svc.TopHeroes()
		if err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, l[:10])
	}
}
