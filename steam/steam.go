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

func (c *Client) Authenticate(r *http.Request) (steamID string, err error) {
	panic("implement me")
}

func (c *Client) Verify(r *http.Request) (*core.SteamPlayer, error) {
	oid := NewOpenId(r)
	mode := oid.Mode()
	if mode == "cancel" {
		return nil, fmt.Errorf("authorization cancelled")
	}

	player, err := oid.ValidateAndGetUser(c.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("could not validate user: %s", err)
	}

	return &core.SteamPlayer{
		ID:     player.SteamId,
		Name:   player.PersonaName,
		URL:    player.ProfileUrl,
		Avatar: player.AvatarFull,
	}, nil
}

func (c *Client) Player(steamID string) (core.SteamPlayer, error) {
	panic("implement me")
}
