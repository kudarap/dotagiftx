package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kudarap/dotagiftx"
)

const userCacheExpr = time.Minute * 5

func handleProfile(svc dotagiftx.UserService, cache cacheManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		cacheKey, noCache := cacheKeyFromRequest(r)
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

func handlePublicProfile(svc dotagiftx.UserService, cache cacheManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		cacheKey, noCache := cacheKeyFromRequest(r)
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

func handleProcSubscription(svc dotagiftx.UserService, cache cacheManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := struct {
			SubscriptionID string `json:"subscription_id"`
		}{}
		if err := parseForm(r, &form); err != nil {
			respondError(w, err)
			return
		}

		u, err := svc.ProcessSubscription(r.Context(), form.SubscriptionID)
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

func handleBlacklisted(svc dotagiftx.UserService, cache cacheManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		cacheKey, noCache := cacheKeyFromRequest(r)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		opts, err := findOptsFromURL(r.URL, &dotagiftx.Item{})
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
			list = []dotagiftx.User{}
		}

		go cache.Set(cacheKey, list, time.Hour*24)

		respondOK(w, list)
	}
}

const userVanityCacheExpr = time.Hour

type vanityUserResp struct {
	dotagiftx.User

	IsRegistered  bool      `json:"is_registered"`
	SteamAvatar   string    `json:"steam_avatar"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
}

// TODO this should be place on service
func handleVanityProfile(svc dotagiftx.UserService, steam dotagiftx.SteamClient, cache cacheManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		cacheKey, noCache := cacheKeyFromRequest(r)
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

func handleUserSubscriptionWebhook(svc dotagiftx.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := svc.UpdateSubscriptionFromWebhook(r.Context(), r); err != nil {
			respondError(w, err)
			return
		}
		respondOK(w, nil)
	}
}

func handleUserManualSubscription(svc dotagiftx.UserService, cache cacheManager, divineKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := isValidDivineKey(r, divineKey); err != nil {
			respondError(w, err)
			return
		}

		var form dotagiftx.ManualSubscriptionParam
		if err := parseForm(r, &form); err != nil {
			respondError(w, err)
			return
		}

		u, err := svc.ProcessManualSubscription(r.Context(), form)
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
