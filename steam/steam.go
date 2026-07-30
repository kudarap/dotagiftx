package steam

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/kudarap/dotagiftx"
)

const (
	VanityURLPrefix  = "https://steamcommunity.com/id/"
	profileURLPrefix = "https://steamcommunity.com/profiles/"

	vanityCacheExpr = time.Hour * 24
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
	cache  cacheReadWriter
}

// New create new steam client instance.
func New(c Config, ca cacheReadWriter) (*Client, error) {
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
	url, ok := cleanProfileURL(rawURL)
	if !ok {
		return "", fmt.Errorf("could not resolve profile URL: %s", rawURL)
	}

	// SteamID might be present on the URL provided.
	if after, ok := strings.CutPrefix(url, profileURLPrefix); ok {
		return after, nil
	}

	// URL provided could be a vanity type and need to be resolved to
	// get steam id.
	vanity := strings.TrimPrefix(url, VanityURLPrefix)
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

type cacheReadWriter interface {
	Set(key string, val any, expr time.Duration) error
	Get(key string) (val string, err error)
}

var reSteamID = regexp.MustCompile("[0-9]{17}")

func cleanProfileURL(rawURL string) (url string, ok bool) {
	// vanity url mode
	if v, ok := strings.CutPrefix(rawURL, VanityURLPrefix); ok {
		s := strings.Split(v, "/")
		if len(s) == 0 {
			return "", false
		}
		return fmt.Sprintf("%s%s", VanityURLPrefix, s[0]), true
	}

	id := reSteamID.FindString(rawURL)
	if id == "" {
		return "", false
	}
	return fmt.Sprintf("%s%s", profileURLPrefix, id), true
}
