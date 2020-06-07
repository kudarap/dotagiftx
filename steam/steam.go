package steam

import (
	"fmt"
	"net/http"

	"github.com/kudarap/dota2giftables/core"
)

// Config represents steam config.
type Config struct {
	Key string
}

// Client represents steam client.
type Client struct {
	ApiKey string
}

// New create new steam client instance.
func New(c Config) (*Client, error) {
	return &Client{c.Key}, nil
}

func (c *Client) AuthorizeURL(r *http.Request) (redirectURL string, err error) {
	oid := NewOpenId(r)
	if oid.Mode() != "" {
		err = fmt.Errorf("could not get redirect URL: %s", oid.Mode())
		return
	}

	return oid.AuthUrl(), nil
}

func (c *Client) Authenticate(r *http.Request) (*core.SteamPlayer, error) {
	oid := NewOpenId(r)
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
	su, err := GetPlayerSummaries(steamID, c.ApiKey)
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
