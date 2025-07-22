package steam

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/kudarap/dotagiftx"
)

// Vanity URL prefixes.
const (
	VanityPrefixID      = "https://steamcommunity.com/id/"
	VanityPrefixProfile = "https://steamcommunity.com/profiles/"
	vanityCacheExpr     = time.Hour * 24
)

// Config represents steam config.
type Config struct {
	Key    string
	Realm  string
	Return string
}

// Client represents a steam client.
type Client struct {
	config Config
	cache  cache
}

// New create new steam client instance.
func New(c Config, ca cache) (*Client, error) {
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

func (c *Client) Authenticate(r *http.Request) (*dotagiftx.SteamPlayer, error) {
	oid := NewOpenId(r, c.config)
	m := oid.Mode()
	if m == "cancel" {
		return nil, fmt.Errorf("authorization cancelled")
	}

	id, err := oid.ValidateAndGetId()
	if err != nil {
		return nil, fmt.Errorf("could not validate player: %s", err)
	}

	p, err := c.Player(id)
	if err != nil {
		return nil, fmt.Errorf("could not get player: %s", err)
	}
	return p, nil
}

func (c *Client) Player(steamID string) (*dotagiftx.SteamPlayer, error) {
	su, err := GetPlayerSummaries(steamID, c.config.Key)
	if err != nil {
		return nil, fmt.Errorf("could not get player: %s", err)
	}

	return &dotagiftx.SteamPlayer{
		ID:     su.SteamId,
		Name:   su.PersonaName,
		URL:    su.ProfileUrl,
		Avatar: su.AvatarFull,
	}, nil
}

func (c *Client) ResolveVanityURL(rawURL string) (steamID string, err error) {
	rawURL = strings.TrimRight(rawURL, "/")

	// SteamID might be present on the URL provided.
	if strings.HasPrefix(rawURL, VanityPrefixProfile) {
		return strings.TrimPrefix(rawURL, VanityPrefixProfile), nil
	}

	vanity := strings.TrimPrefix(rawURL, VanityPrefixID)
	cacheKey := fmt.Sprintf("steam/resolvedvanity:%s", vanity)
	if hit, _ := c.cache.Get(cacheKey); hit != "" {
		return strings.ReplaceAll(hit, `"`, ""), nil
	}

	steamID, err = ResolveVanityURL(vanity, c.config.Key)
	if err != nil {
		return
	}

	err = c.cache.Set(cacheKey, steamID, vanityCacheExpr)
	return
}

type cache interface {
	Set(key string, val interface{}, expr time.Duration) error
	Get(key string) (val string, err error)
}
