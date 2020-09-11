package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/kudarap/dotagiftx/core"
)

func handleProfile(svc core.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := svc.UserFromContext(r.Context())
		if err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, u)
	}
}

func handlePublicProfile(svc core.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		u, err := svc.User(id)
		if err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, u)
	}
}
