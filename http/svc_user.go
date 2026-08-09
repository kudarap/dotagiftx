package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	dotagiftx2 "github.com/kudarap/dotagiftx/dotagiftx"
	"github.com/kudarap/dotagiftx/steam"
)

const userCacheExpr = time.Minute * 5

// userService provides access to user service methods used by http handlers.
type userService interface {
	// FlaggedUsers returns a list of flagged/reported users.
	FlaggedUsers(ctx context.Context, opts dotagiftx2.FindOpts) ([]dotagiftx2.User, error)

	// User returns user details by id.
	User(ctx context.Context, id string) (*dotagiftx2.User, error)

	// UserFromContext returns user details from context.
	UserFromContext(context.Context) (*dotagiftx2.User, error)

	// CreateSubscription creates a subscription for the current user.
	CreateSubscription(ctx context.Context, planID string) (subscriptionID string, err error)

	// ProcessSubscription validates and processes subscription features.
	ProcessSubscription(ctx context.Context, subscriptionID string) (*dotagiftx2.User, error)

	// UpdateSubscriptionFromWebhook handles user subscription updates from http request.
	UpdateSubscriptionFromWebhook(ctx context.Context, r *http.Request) (*dotagiftx2.User, error)

	// ProcessManualSubscription processes manual subscription.
	ProcessManualSubscription(ctx context.Context, form dotagiftx2.ManualSubscriptionParam) (*dotagiftx2.User, error)
}

// steamClient provides access to steam API methods used by http handlers.
type steamClient interface {
	// Player returns player summary base on steamID.
	Player(steamID string) (*dotagiftx2.SteamPlayer, error)

	// ResolveVanityURL returns steam id from profile url.
	ResolveVanityURL(url string) (steamID string, err error)
}

func handleProfile(svc userService, cache cacheManager, logger *slog.Logger) http.HandlerFunc {
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

		go func() {
			if err := cache.Set(cacheKey, u, userCacheExpr); err != nil {
				logger.ErrorContext(r.Context(), "could not save cache on profile", "error", err)
			}
		}()

		respondOK(w, u)
	}
}

func handlePublicProfile(svc userService, cache cacheManager, logger *slog.Logger) http.HandlerFunc {
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
		u, err := svc.User(r.Context(), id)
		if err != nil {
			respondError(w, err)
			return
		}

		go func() {
			if err := cache.Set(cacheKey, u, userCacheExpr); err != nil {
				logger.ErrorContext(r.Context(), "could not save cache on public profile", "error", err)
			}
		}()

		respondOK(w, u)
	}
}

func handleBlacklisted(svc userService, cache cacheManager, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for cache hit and render them.
		cacheKey, noCache := cacheKeyFromRequest(r)
		if !noCache {
			if hit, _ := cache.Get(cacheKey); hit != "" {
				respondOK(w, hit)
				return
			}
		}

		opts, err := findOptsFromURL(r.URL, &dotagiftx2.Item{})
		if err != nil {
			respondError(w, err)
			return
		}
		list, err := svc.FlaggedUsers(r.Context(), opts)
		if err != nil {
			respondError(w, err)
			return
		}
		if list == nil {
			list = []dotagiftx2.User{}
		}

		go func() {
			if err := cache.Set(cacheKey, list, time.Hour*24); err != nil {
				logger.ErrorContext(r.Context(), "could not save cache on blacklists", "error", err)
			}
		}()

		respondOK(w, list)
	}
}

const userVanityCacheExpr = time.Hour

type vanityUserResp struct {
	dotagiftx2.User

	IsRegistered  bool      `json:"is_registered"`
	SteamAvatar   string    `json:"steam_avatar"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
}

// TODO this should be place on service
func handleVanityProfile(svc userService, steamClient steamClient, cache cacheManager, logger *slog.Logger) http.HandlerFunc {
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
		steamID, err := steamClient.ResolveVanityURL(fmt.Sprintf(steam.VanityURLPrefix+"%s", id))
		if err != nil {
			respondError(w, err)
			return
		}
		vUser.SteamID = steamID

		// Get user data if its registered.
		u, _ := svc.User(r.Context(), steamID)
		if u != nil {
			vUser.User = *u
			vUser.IsRegistered = true
		} else {
			// Otherwise, get it from steam API.
			sp, err := steamClient.Player(steamID)
			if err != nil {
				respondError(w, err)
				return
			}
			vUser.Name = sp.Name
			vUser.URL = sp.URL
			vUser.SteamAvatar = sp.Avatar
		}

		vUser.LastUpdatedAt = time.Now()

		go func() {
			if err := cache.Set(cacheKey, vUser, userVanityCacheExpr); err != nil {
				logger.ErrorContext(r.Context(), "could not save cache on vanity profile", "error", err)
			}
		}()
		respondOK(w, vUser)
	}
}

func handleCreateSubscription(svc userService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := struct {
			PlanID string `json:"plan_id"`
		}{}
		if err := parseForm(r, &form); err != nil {
			respondError(w, err)
			return
		}

		subID, err := svc.CreateSubscription(r.Context(), form.PlanID)
		if err != nil {
			respondError(w, err)
			return
		}

		respondOK(w, struct {
			ID string `json:"id"`
		}{subID})
	}
}

func handleProcSubscription(svc userService, cache cacheManager, logger *slog.Logger) http.HandlerFunc {
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
			if err := cache.BulkDel(fmt.Sprintf("users/%s*", u.SteamID)); err != nil {
				logger.ErrorContext(r.Context(), "could not invalidate user cache", "error", err)
			}
			if err := cache.BulkDel(marketCacheKeyPrefix); err != nil {
				logger.ErrorContext(r.Context(), "could not invalidate market cache", "error", err)
			}
		}()
		respondOK(w, u)
	}
}

func handleUserSubscriptionWebhook(svc userService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := svc.UpdateSubscriptionFromWebhook(r.Context(), r); err != nil {
			respondError(w, err)
			return
		}
		respondOK(w, nil)
	}
}

func handleUserManualSubscription(svc userService, cache cacheManager, divineKey string, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := isValidDivineKey(r, divineKey); err != nil {
			respondError(w, err)
			return
		}

		var form dotagiftx2.ManualSubscriptionParam
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
			if err := cache.BulkDel(fmt.Sprintf("users/%s*", u.SteamID)); err != nil {
				logger.ErrorContext(r.Context(), "could not invalidate user cache", "error", err)
			}
			if err := cache.BulkDel(marketCacheKeyPrefix); err != nil {
				logger.ErrorContext(r.Context(), "could not invalidate market cache", "error", err)
			}
		}()
		respondOK(w, u)
	}
}
