package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kudarap/dotagiftx/cache"
	"github.com/kudarap/dotagiftx/phantasm"
)

func handlePhantasmWebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "steam_id")
		key := fmt.Sprintf("phantasm_%s", id)
		ttl := time.Hour

		var inventory phantasm.Inventory
		if err := parseForm(r, &inventory); err != nil {
			respondError(w, err)
			return
		}
		if err := cache.Set(key, inventory, ttl); err != nil {
			respondError(w, err)
			return
		}
		respondOK(w, "ok")
	}
}
