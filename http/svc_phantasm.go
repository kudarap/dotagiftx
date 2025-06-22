package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kudarap/dotagiftx/phantasm"
)

func handlePhantasmWebhook(svc *phantasm.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "steam_id")
		if err := svc.SaveInventory(r.Context(), id, r.Body); err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, "ok")
	}
}
