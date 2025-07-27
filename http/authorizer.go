package http

import (
	"net/http"

	"github.com/kudarap/dotagiftx"
)

func (s *Server) authorizer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate token from header.
		c, err := ParseFromHeader(r.Header)
		if err != nil {
			respondError(w, dotagiftx.AuthErrNoAccess.X(err))
			return
		}

		// Checks auth level required.
		// Inject auth details to context that will later be used as
		// context user and authorizer level.
		ctx := dotagiftx.AuthToContext(r.Context(), &dotagiftx.Auth{
			UserID: c.UserID,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
