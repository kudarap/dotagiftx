package http

import (
	"context"
	"log/slog"
	"net/http"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
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

func logger(logger *slog.Logger) func(next http.Handler) http.Handler {
	logFormat := httplog.SchemaECS.Concise(true)
	return httplog.RequestLogger(logger, &httplog.Options{
		// Level defines the verbosity of the request logs:
		// slog.LevelDebug - log all responses (incl. OPTIONS)
		// slog.LevelInfo  - log responses (excl. OPTIONS)
		// slog.LevelWarn  - log 4xx and 5xx responses only (except for 429)
		// slog.LevelError - log 5xx responses only
		Level: slog.LevelInfo,

		// Set log output to Elastic Common Schema (ECS) format.
		Schema: logFormat,

		// RecoverPanics recovers from panics occurring in the underlying HTTP handlers
		// and middlewares. It returns HTTP 500 unless response status was already set.
		//
		// NOTE: Panics are logged as errors automatically, regardless of this setting.
		RecoverPanics: true,

		// Optionally, filter out some request logs.
		Skip: func(r *http.Request, status int) bool {
			if isHealthChech(r) {
				return true
			}

			return status == 404 || status == 405
		},

		// Optionally, log selected request/response headers explicitly.
		LogRequestHeaders:  []string{"Origin", requestIDHeader},
		LogResponseHeaders: []string{},

		// Optionally, enable logging of request/response body based on custom conditions.
		// Useful for debugging payload issues in development.
		LogRequestBody:  isDebugHeaderSet,
		LogResponseBody: isDebugHeaderSet,
	})
}

func isDebugHeaderSet(r *http.Request) bool {
	return r.Header.Get("Debug") == "true"
}

func isHealthChech(r *http.Request) bool {
	return r.URL.Path == "/" && r.URL.Query().Has("health")
}
