package http

import (
	"context"
	"net/http"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

const (
	requestIDKey          = chimiddleware.RequestIDKey
	requestIDHeader       = "X-Request-Id"
	vercelIDRequestHeader = "X-Vercel-Id"
)

func vercelRequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := r.Header.Get(vercelIDRequestHeader)
		if requestID != "" {
			ctx = context.WithValue(ctx, requestIDKey, requestID)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// cors middleware injects cors headers to each request.
func cors(next http.Handler) http.Handler {
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

func requestIDWriter(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		rid, _ := r.Context().Value(requestIDKey).(string)
		w.Header().Set(requestIDHeader, rid)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
