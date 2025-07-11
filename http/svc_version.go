package http

import (
	"net/http"

	"github.com/kudarap/dotagiftx"
)

func handleInfo(v *dotagiftx.Version) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respondOK(w, v)
	}
}
