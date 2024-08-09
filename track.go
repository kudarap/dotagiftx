package dgx

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Track error types.
const (
	TrackErrNotFound Errors = iota + 4000
)

// sets error text definition.
func init() {
	appErrorText[TrackErrNotFound] = "track details not found"
}

// Track types.
const (
	TrackTypeView        = "v"
	TrackTypeSearch      = "s"
	TrackTypeProfileView = "p"
	//TrackTypeItemView           = 300
	//TrackTypeCatalogSearch      = 310
	//TrackTypeProfileClick       = 110
	//TrackTypeProfileView        = 100
	//TrackTypeMarketListed       = 220
	//TrackTypeMarketReserved     = 230
	//TrackTypeMarketSold         = 240
	//TrackTypeMarketBidCompleted = 241
	//TrackTypeMarketRemoved      = 250
	//TrackTypeMarketCancelled    = 260
	//TrackTypeMarketExpired      = 270
)

type (
	// Track represents tracking data.
	Track struct {
		ID         string     `json:"id"           db:"id,omitempty"`
		Type       string     `json:"type"         db:"type,omitempty,indexed"`
		ItemID     string     `json:"item_id"      db:"item_id,omitempty,indexed"`
		UserID     string     `json:"user_id"      db:"user_id,omitempty,indexed"`
		Keyword    string     `json:"keyword"      db:"keyword,omitempty"`
		ClientIP   string     `json:"client_ip"    db:"client_ip,omitempty"`
		UserAgent  string     `json:"user_agent"   db:"user_agent,omitempty"`
		Referer    string     `json:"referer"      db:"referer,omitempty"`
		Cookies    []string   `json:"cookies"      db:"cookies,omitempty"`
		SessUserID string     `json:"sess_user_id" db:"sess_user_id,omitempty"`
		CreatedAt  *time.Time `json:"created_at"   db:"created_at,omitempty,indexed"`
		UpdatedAt  *time.Time `json:"updated_at"   db:"updated_at,omitempty"`
	}

	// TrackService provides access to track service.
	TrackService interface {
		// Tracks returns a list of tracks.
		Tracks(FindOpts) ([]Track, *FindMetadata, error)

		// Track returns track details by id.
		Track(id string) (*Track, error)

		// CreateFromRequest saves new track from http request. Primarily used on client side.
		CreateFromRequest(r *http.Request) error

		// CreateSearchKeyword saves new keyword tracking data.
		CreateSearchKeyword(r *http.Request, keyword string) error
	}

	// TrackStorage defines operation for track records.
	TrackStorage interface {
		// Find returns a list of tracks from data store.
		Find(FindOpts) ([]Track, error)

		// Count returns number of tracks from data store.
		Count(FindOpts) (int, error)

		// Get returns track details by id from data store.
		Get(id string) (*Track, error)

		// Create persists a new track to data store.
		Create(*Track) error

		// ThisWeekKeywords returns top search keywords this week.
		TopKeywords() ([]SearchKeywordScore, error)
	}
)

const (
	trackTypeKey    = "t"
	trackItemIDKey  = "i"
	trackUserIDKey  = "u"
	trackKeywordKey = "k"
)

const authCookieName = "dgAu"

// SetDefaults sets default values from http.Request.
func (t *Track) SetDefaults(r *http.Request) {
	q := r.URL.Query()
	t.Type = q.Get(trackTypeKey)
	t.ItemID = q.Get(trackItemIDKey)
	t.UserID = q.Get(trackUserIDKey)
	t.Keyword = q.Get(trackKeywordKey)
	t.Referer = r.Referer()
	t.UserAgent = r.UserAgent()

	ip, err := userIPFromRequest(r)
	if err != nil {
		t.ClientIP = r.RemoteAddr
	} else {
		t.ClientIP = ip.String()
	}

	var sessCookie string
	for _, cookie := range r.Cookies() {
		s, _ := url.PathUnescape(cookie.String())
		t.Cookies = append(t.Cookies, s)

		if cookie.Name == authCookieName {
			sessCookie = s
		}
	}

	// Extract session user id from cookie.
	var au Auth
	sessCookie = strings.TrimPrefix(sessCookie, authCookieName+"=")
	_ = json.Unmarshal([]byte(sessCookie), &au)
	t.SessUserID = au.UserID
}

func userIPFromRequest(req *http.Request) (net.IP, error) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}
	return userIP, nil
}
