package http

import (
	"fmt"
	"net/http"

	"github.com/kudarap/dotagiftx"
)

func handleHammerBan(svc dotagiftx.HammerService, cache Cache) http.HandlerFunc {
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

		go resetProfileListingCache(u.SteamID, cache)
		respondOK(w, u)
	}
}

func handleHammerSuspend(svc dotagiftx.HammerService, cache Cache) http.HandlerFunc {
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

		go resetProfileListingCache(u.SteamID, cache)
		respondOK(w, u)
	}
}

func handleHammerLift(svc dotagiftx.HammerService, cache Cache) http.HandlerFunc {
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

func resetProfileListingCache(steamID string, cache Cache) {
	cache.BulkDel("blacklists")
	cache.BulkDel(fmt.Sprintf("users/%s*", steamID))
	cache.BulkDel(marketCacheKeyPrefix)
}
