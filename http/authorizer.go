package http

import (
	"net/http"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/errors"
	"github.com/kudarap/dotagiftx/gokit/http/jwt"
)

func (s *Server) authorizer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate token from header.
		c, err := jwt.ParseFromHeader(r.Header)
		if err != nil {
			respondError(w, errors.New(core.AuthErrNoAccess, err))
			return
		}

		// Checks auth level required.

		// Inject auth details to context that will later be use as
		// context user and authorizer level.
		ctx := core.AuthToContext(r.Context(), &core.Auth{
			UserID: c.UserID,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
