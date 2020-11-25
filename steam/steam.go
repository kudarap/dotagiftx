package steam

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/kudarap/dotagiftx/core"
)

// Config represents steam config.
type Config struct {
	Key    string
	Realm  string
	Return string
}

// Client represents steam client.
type Client struct {
	config Config
}

// New create new steam client instance.
func New(c Config) (*Client, error) {
	return &Client{c}, nil
}

func (c *Client) AuthorizeURL(r *http.Request) (redirectURL string, err error) {
	// Check callback URL override and create a config.
	cb := r.URL.Query().Get("callback")
	if cb != "" {
		u, _ := url.Parse(cb)
		c.config.Realm = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
		c.config.Return = cb
	}

	oid := NewOpenId(r, c.config)
	if oid.Mode() != "" {
		err = fmt.Errorf("could not get redirect URL: %s", oid.Mode())
		return
	}

	return oid.AuthUrl(), nil
}

func (c *Client) Authenticate(r *http.Request) (*core.SteamPlayer, error) {
	oid := NewOpenId(r, c.config)
	m := oid.Mode()
	if m == "cancel" {
		return nil, fmt.Errorf("authorization cancelled")
	}

	id, err := oid.ValidateAndGetId()
	if err != nil {
		return nil, fmt.Errorf("could not validate player: %s", err)
	}

	return c.Player(id)
}

func (c *Client) Player(steamID string) (*core.SteamPlayer, error) {
	su, err := GetPlayerSummaries(steamID, c.config.Key)
	if err != nil {
		return nil, fmt.Errorf("could not get player: %s", err)
	}

	return &core.SteamPlayer{
		ID:     su.SteamId,
		Name:   su.PersonaName,
		URL:    su.ProfileUrl,
		Avatar: su.AvatarFull,
	}, nil
}
