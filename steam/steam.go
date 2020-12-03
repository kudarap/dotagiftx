package steam

import (
	"fmt"
	"net/http"
	"strings"
	"time"

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
	cache  core.Cache
}

// New create new steam client instance.
func New(c Config, ca core.Cache) (*Client, error) {
	return &Client{c, ca}, nil
}

func (c *Client) AuthorizeURL(r *http.Request) (redirectURL string, err error) {
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

// Vanity URL prefixes.
const (
	VanityPrefixID      = "https://steamcommunity.com/id/"
	VanityPrefixProfile = "https://steamcommunity.com/profiles/"
)

func (c *Client) ResolveVanityURL(rawURL string) (steamID string, err error) {
	cacheKey := fmt.Sprintf("steam/vanity/%s", rawURL)
	if hit, err := c.cache.Get(cacheKey); err != nil {
		fmt.Println("cache hit!")
		return hit, err
	}

	rawURL = strings.TrimRight(rawURL, "/")

	// SteamID might be present on the URL already.
	if strings.HasPrefix(rawURL, VanityPrefixProfile) {
		return strings.TrimPrefix(rawURL, VanityPrefixProfile), nil
	}

	// Its probably steam ID.
	if !strings.HasPrefix(rawURL, VanityPrefixID) {
		err = fmt.Errorf("could not parse URL (%s)", rawURL)
		return
	}

	v := strings.TrimPrefix(rawURL, VanityPrefixID)
	steamID, err = ResolveVanityURL(v, c.config.Key)
	if err != nil {
		return
	}

	c.cache.Set(cacheKey, steamID, time.Hour*24)
	return
}
