package middleware

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
)

// RecovererWithLog writes to logs when internal error occurred.
func RecovererWithLog(logWriter io.Writer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					rqd, err := httputil.DumpRequest(r, true)
					if err != nil {
						http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
						return
					}

					fmt.Fprintf(logWriter, "Panic: %+v\n\nReq Dump: %sStack Trace: %s", rvr, rqd, debug.Stack())
					//debug.PrintStack()

					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// Recoverer writes to logs os standard error.
func Recoverer(next http.Handler) http.Handler {
	return RecovererWithLog(os.Stderr)(next)
}
