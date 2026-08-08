package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/kudarap/dotagiftx"
)

func handleHammerBan(svc hammerService, cache cacheManager, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p dotagiftx.HammerParams
		if err := parseForm(r, &p); err != nil {
			respondError(w, err)
			return
		}

		u, err := svc.Ban(r.Context(), p)
		if err != nil {
			respondError(w, err)
			return
		}

		go resetProfileListingCache(context.WithoutCancel(r.Context()), u.SteamID, cache, logger)
		respondOK(w, u)
	}
}

func handleHammerSuspend(svc hammerService, cache cacheManager, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p dotagiftx.HammerParams
		if err := parseForm(r, &p); err != nil {
			respondError(w, err)
			return
		}

		u, err := svc.Suspend(r.Context(), p)
		if err != nil {
			respondError(w, err)
			return
		}

		go resetProfileListingCache(context.WithoutCancel(r.Context()), u.SteamID, cache, logger)
		respondOK(w, u)
	}
}

func handleHammerLift(svc hammerService, cache cacheManager, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := struct {
			SteamID         string `json:"steam_id"`
			RestoreListings bool   `json:"restore_listings"`
		}{}
		if err := parseForm(r, &p); err != nil {
			respondError(w, err)
			return
		}

		if err := svc.Lift(r.Context(), p.SteamID, p.RestoreListings); err != nil {
			respondError(w, err)
			return
		}

		go resetProfileListingCache(context.WithoutCancel(r.Context()), p.SteamID, cache, logger)
		respondOK(w, newMsg("hammer lifted"))
	}
}

func resetProfileListingCache(ctx context.Context, steamID string, cache cacheManager, logger *slog.Logger) {
	if err := cache.BulkDel("blacklists"); err != nil {
		logger.ErrorContext(ctx, "could not invalidate blacklists cache", "error", err)
	}
	if err := cache.BulkDel(fmt.Sprintf("users/%s*", steamID)); err != nil {
		logger.ErrorContext(ctx, "could not invalidate user cache", "error", err)
	}
	if err := cache.BulkDel(marketCacheKeyPrefix); err != nil {
		logger.ErrorContext(ctx, "could not invalidate market cache", "error", err)
	}
}
