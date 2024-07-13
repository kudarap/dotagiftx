package http

import (
	"net/http"

	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/errors"
	"github.com/kudarap/dotagiftx/gokit/http/jwt"
)

func (s *Server) authorizer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate token from header.
		c, err := jwt.ParseFromHeader(r.Header)
		if err != nil {
			respondError(w, errors.New(dgx.AuthErrNoAccess, err))
			return
		}

		// Checks auth level required.

		// Inject auth details to context that will later be use as
		// context user and authorizer level.
		ctx := dgx.AuthToContext(r.Context(), &dgx.Auth{
			UserID: c.UserID,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
