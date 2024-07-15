package steam

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	dgx "github.com/kudarap/dotagiftx"
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
	cache  dgx.Cache
}

// New create new steam client instance.
func New(c Config, ca dgx.Cache) (*Client, error) {
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

func (c *Client) Authenticate(r *http.Request) (*dgx.SteamPlayer, error) {
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

func (c *Client) Player(steamID string) (*dgx.SteamPlayer, error) {
	su, err := GetPlayerSummaries(steamID, c.config.Key)
	if err != nil {
		return nil, fmt.Errorf("could not get player: %s", err)
	}

	return &dgx.SteamPlayer{
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
	vanityCacheExpr     = time.Hour * 24
)

func (c *Client) ResolveVanityURL(rawURL string) (steamID string, err error) {
	rawURL = strings.TrimRight(rawURL, "/")

	// SteamID might be present on the URL provided.
	if strings.HasPrefix(rawURL, VanityPrefixProfile) {
		return strings.TrimPrefix(rawURL, VanityPrefixProfile), nil
	}

	// Its probably steam ID.
	//if !strings.HasPrefix(rawURL, VanityPrefixID) {
	//	err = fmt.Errorf("could not parse URL (%s)", rawURL)
	//	return
	//}

	vanity := strings.TrimPrefix(rawURL, VanityPrefixID)
	cacheKey := fmt.Sprintf("steam/resolvedvanity:%s", vanity)
	if hit, _ := c.cache.Get(cacheKey); hit != "" {
		return strings.ReplaceAll(hit, `"`, ""), nil
	}

	steamID, err = ResolveVanityURL(vanity, c.config.Key)
	if err != nil {
		return
	}

	c.cache.Set(cacheKey, steamID, vanityCacheExpr)
	return
}
