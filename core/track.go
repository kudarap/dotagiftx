package core

import (
	"fmt"
	"net"
	"net/http"
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
	TrackTypeView   = "v"
	TrackTypeSearch = "s"
)

type (
	Track struct {
		ID        string     `json:"id"         db:"id,omitempty"`
		Type      string     `json:"type"       db:"type,omitempty"`
		ItemID    string     `json:"item_id"    db:"item_id,omitempty"`
		Keyword   string     `json:"keyword"    db:"keyword,omitempty"`
		ClientIP  string     `json:"client_ip"  db:"client_ip,omitempty"`
		UserAgent string     `json:"user_agent" db:"user_agent,omitempty"`
		Referer   string     `json:"referer"    db:"referer,omitempty"`
		CreatedAt *time.Time `json:"created_at" db:"created_at,omitempty"`
		UpdatedAt *time.Time `json:"updated_at" db:"updated_at,omitempty"`
	}

	// TrackService provides access to track service.
	TrackService interface {
		// Tracks returns a list of tracks.
		Tracks(FindOpts) ([]Track, *FindMetadata, error)

		// Track returns track details by id.
		Track(id string) (*Track, error)

		// Create saves new track.
		Create(r *http.Request) error
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
	}
)

const (
	trackTypeKey    = "t"
	trackItemIDKey  = "i"
	trackKeywordKey = "k"
)

// SetDefaults sets default values from http.Request.
func (t *Track) SetDefaults(r *http.Request) {
	q := r.URL.Query()
	t.Type = q.Get(trackTypeKey)
	t.ItemID = q.Get(trackItemIDKey)
	t.Keyword = q.Get(trackKeywordKey)
	t.Referer = r.Referer()
	t.UserAgent = r.UserAgent()

	ip, err := userIPFromRequest(r)
	if err != nil {
		t.ClientIP = r.RemoteAddr
	} else {
		t.ClientIP = ip.String()
	}
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
