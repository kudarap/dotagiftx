package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"runtime/debug"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

const vercelIDRequestHeader = "X-Vercel-Id"

func vercelRequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := r.Header.Get(vercelIDRequestHeader)
		if requestID != "" {
			ctx = context.WithValue(ctx, chimiddleware.RequestIDKey, requestID)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

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
