package http

import (
	"fmt"
	"net/http"

	"github.com/kudarap/dotagiftx/core"
)

func handleHammerBan(svc core.HammerService, cache core.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p core.HammerParams
		if err := parseForm(r, &p); err != nil {
			respondError(w, err)
			return
		}

		u, err := svc.Ban(r.Context(), p)
		if err != nil {
			respondError(w, err)
			return
		}

		go resetProfileListingCache(u.SteamID, cache)
		respondOK(w, u)
	}
}

func handleHammerSuspend(svc core.HammerService, cache core.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p core.HammerParams
		if err := parseForm(r, &p); err != nil {
			respondError(w, err)
			return
		}

		u, err := svc.Suspend(r.Context(), p)
		if err != nil {
			respondError(w, err)
			return
		}

		go resetProfileListingCache(u.SteamID, cache)
		respondOK(w, u)
	}
}

func handleHammerLift(svc core.HammerService, cache core.Cache) http.HandlerFunc {
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

		go resetProfileListingCache(p.SteamID, cache)
		respondOK(w, newMsg("hammer lifted"))
	}
}

func resetProfileListingCache(steamID string, cache core.Cache) {
	cache.BulkDel(fmt.Sprintf("/blacklists/%s", steamID))
	cache.BulkDel(fmt.Sprintf("/users/%s", steamID))
	cache.BulkDel(marketCacheKeyPrefix)
}
