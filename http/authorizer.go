package http

import (
	"net/http"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/errors"
)

func (s *Server) authorizer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate token from header.
		c, err := ParseFromHeader(r.Header)
		if err != nil {
			respondError(w, errors.New(dotagiftx.AuthErrNoAccess, err))
			return
		}

		// Checks auth level required.

		// Inject auth details to context that will later be use as
		// context user and authorizer level.
		ctx := dotagiftx.AuthToContext(r.Context(), &dotagiftx.Auth{
			UserID: c.UserID,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
