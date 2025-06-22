package http

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kudarap/dotagiftx/phantasm"
)

func handlePhantasmWebhook(svc *phantasm.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "steam_id")

		var b bytes.Buffer
		n, err := io.Copy(&b, r.Body)
		if err != nil {
			respondError(w, err)
			return
		}

		log.Printf("received %d bytes from steam", n)
		if err := svc.SaveInventory(r.Context(), id, &b); err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, "ok")
	}
}
