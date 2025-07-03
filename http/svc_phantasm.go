package http

import (
	"errors"
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

func handlePhantasmCrawl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		args := map[string]interface{}{
			"steam_id": r.URL.Query().Get("steam_id"),
		}
		resp := phantasm.Main(args)

		code, ok := resp["statusCode"].(int)
		if !ok {
			respondError(w, errors.New("invalid status code"))
			return
		}
		respond(w, code, resp["body"])
	}
}
