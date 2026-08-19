package http

import (
	"net/http"

	"github.com/kudarap/dotagiftx/dotagiftx"
)

func handleInfo(v *dotagiftx.Version) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v = v.SetUptime()
		respondOK(w, v)
	}
}
