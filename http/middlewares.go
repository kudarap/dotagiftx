package http

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
)

// CORS middleware injects CORS headers to each request.
func CORS(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "HEAD, GET, POST, PATCH, DELETE, OPTIONS")

		// NOTE handle OPTIONS and HEAD method to respond immediately.
		if r.Method == http.MethodHead || r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

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

func CORSx(allowed []string) func(http.Handler) http.Handler {
	isOriginAllowed := func(origin string) bool {
		if len(allowed) == 0 {
			return true
		}

		for _, aa := range allowed {
			if origin == aa {
				return true
			}
		}

		return false
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" && isOriginAllowed(origin) {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("origin not allowed"))
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "HEAD, GET, POST, PATCH, DELETE, OPTIONS")

			// NOTE handle OPTIONS and HEAD method to respond immediately.
			if r.Method == http.MethodHead || r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
