package core

import "net/http"

type (
	// SteamPlayer represents steam player information.
	SteamPlayer struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		URL    string `json:"url"`
		Avatar string `json:"avatar"`
	}

	// SteamClient provides access to Steam API.
	SteamClient interface {
		// AuthorizeURL returns authorization url to steam open id.
		AuthorizeURL(r *http.Request) (redirectURL string, err error)

		// Authenticate returns steamID on valid authorization.
		Authenticate(r *http.Request) (steamID string, err error)

		// Player returns player summary base on steamID.
		Player(steamID string) (SteamPlayer, error)
	}
)
