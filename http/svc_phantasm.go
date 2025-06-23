package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kudarap/dotagiftx/phantasm"
)

func handlePhantasmWebhook(svc *phantasm.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "steam_id")
		secret := r.Header.Get(phantasm.WebhookAuthHeader)
		if err := svc.SaveInventory(r.Context(), id, secret, r.Body); err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, "ok")
	}
}
