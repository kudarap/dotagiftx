package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/kudarap/dotagiftx/core"
)

const userCacheExpr = time.Minute * 5

func handleProfile(svc core.UserService, cache core.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		cacheKey, noCache := core.CacheKeyFromRequest(r)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		u, err := svc.UserFromContext(r.Context())
		if err != nil {
			respondError(w, err)
			return
		}

		go cache.Set(cacheKey, u, userCacheExpr)

		respondOK(w, u)
	}
}

func handlePublicProfile(svc core.UserService, cache core.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		cacheKey, noCache := core.CacheKeyFromRequest(r)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		id := chi.URLParam(r, "id")
		u, err := svc.User(id)
		if err != nil {
			respondError(w, err)
			return
		}

		go cache.Set(cacheKey, u, userCacheExpr)

		respondOK(w, u)
	}
}

const userVanityCacheExpr = time.Hour

type vanityUserResp struct {
	core.User

	IsRegistered  bool      `json:"is_registered"`
	SteamImage    string    `json:"steam_image"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
}

// TODO this should be place on service
func handleVanityProfile(svc core.UserService, steam core.SteamClient, cache core.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		cacheKey, noCache := core.CacheKeyFromRequest(r)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		vUser := new(vanityUserResp)

		// Try to resolve the vanity URL or vanity.
		id := chi.URLParam(r, "id")
		steamID, err := steam.ResolveVanityURL(id)
		if err != nil {
			respondError(w, err)
			return
		}
		vUser.SteamID = steamID

		// Get user data if its registered.
		u, _ := svc.User(steamID)
		if u != nil {
			vUser.User = *u
			vUser.IsRegistered = true
		}

		// Otherwise, get it from steam API.
		sp, err := steam.Player(steamID)
		if err != nil {
			respondError(w, err)
			return
		}
		vUser.Name = sp.Name
		vUser.URL = sp.URL
		vUser.SteamImage = sp.Avatar
		vUser.LastUpdatedAt = time.Now()

		go cache.Set(cacheKey, vUser, userVanityCacheExpr)
		respondOK(w, vUser)
	}
}
