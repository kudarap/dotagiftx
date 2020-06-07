package core

import "net/http"

type (
	// SteamPlayer represents steam player information.
	SteamPlayer struct {
		ID     string `json:"id"     db:"id"`
		Name   string `json:"name"   db:"name"`
		URL    string `json:"url"    db:"url"`
		Avatar string `json:"avatar" db:"avatar"`
	}

	// SteamClient provides access to Steam API.
	SteamClient interface {
		// AuthorizeURL returns authorization url to steam open id.
		AuthorizeURL(r *http.Request) (redirectURL string, err error)

		// Authenticate returns player info on valid authorization.
		Authenticate(r *http.Request) (*SteamPlayer, error)

		// Player returns player summary base on steamID.
		Player(steamID string) (*SteamPlayer, error)
	}
)
