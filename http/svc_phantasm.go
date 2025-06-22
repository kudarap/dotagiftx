package http

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kudarap/dotagiftx/cache"
)

func handlePhantasmWebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "steam_id")
		key := fmt.Sprintf("phantasm_%s", id)
		ttl := time.Hour

		resBody, err := io.ReadAll(r.Body)
		if err != nil {
			respondError(w, err)
			return
		}
		if err = cache.Set(key, resBody, ttl); err != nil {
			respondError(w, err)
			return
		}
		respondOK(w, "ok")
	}
}
