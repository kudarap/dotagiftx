package http

import (
	"fmt"
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

func handleProcSubscription(svc core.UserService, cache core.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := struct {
			SubscriptionID string `json:"subscription_id"`
		}{}
		if err := parseForm(r, &form); err != nil {
			respondError(w, err)
			return
		}

		u, err := svc.ProcSubscription(r.Context(), form.SubscriptionID)
		if err != nil {
			respondError(w, err)
			return
		}

		go func() {
			cache.BulkDel(fmt.Sprintf("users/%s*", u.SteamID))
			cache.BulkDel(marketCacheKeyPrefix)
		}()
		respondOK(w, u)
	}
}

func handleBlacklisted(svc core.UserService, cache core.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		cacheKey, noCache := core.CacheKeyFromRequest(r)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		opts, err := findOptsFromURL(r.URL, &core.Item{})
		if err != nil {
			respondError(w, err)
			return
		}
		list, err := svc.FlaggedUsers(opts)
		if err != nil {
			respondError(w, err)
			return
		}
		if list == nil {
			list = []core.User{}
		}

		go cache.Set(cacheKey, list, userCacheExpr)

		respondOK(w, list)
	}
}

const userVanityCacheExpr = time.Hour

type vanityUserResp struct {
	core.User

	IsRegistered  bool      `json:"is_registered"`
	SteamAvatar   string    `json:"steam_avatar"`
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
		} else {
			// Otherwise, get it from steam API.
			sp, err := steam.Player(steamID)
			if err != nil {
				respondError(w, err)
				return
			}
			vUser.Name = sp.Name
			vUser.URL = sp.URL
			vUser.SteamAvatar = sp.Avatar
		}

		vUser.LastUpdatedAt = time.Now()

		go cache.Set(cacheKey, vUser, userVanityCacheExpr)
		respondOK(w, vUser)
	}
}
