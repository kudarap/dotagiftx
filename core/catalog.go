package core

import (
	"time"
)

// Catalog error types.
const (
	CatalogErrNotFound Errors = iota + 2200
	CatalogErrRequiredID
	CatalogErrIndexing
)

// sets error text definition.
func init() {
	appErrorText[CatalogErrNotFound] = "catalog not found"
	appErrorText[CatalogErrRequiredID] = "catalog id is required"
	appErrorText[CatalogErrIndexing] = "catalog indexing error"
}

type (
	// Catalog represents item market information.
	Catalog struct {
		ID           string   `json:"id"           db:"id,omitempty"`
		Slug         string   `json:"slug"         db:"slug,omitempty,indexed"`
		Name         string   `json:"name"         db:"name,omitempty,indexed"`
		Hero         string   `json:"hero"         db:"hero,omitempty,indexed"`
		Image        string   `json:"image"        db:"image,omitempty"`
		Origin       string   `json:"origin"       db:"origin,omitempty,indexed"`
		Rarity       string   `json:"rarity"       db:"rarity,omitempty,indexed"`
		Contributors []string `json:"-"            db:"contributors,omitempty"`
		ViewCount    int      `json:"view_count"   db:"view_count,omitempty,indexed"`
		// Market summary details.
		Quantity   int        `json:"quantity"    db:"quantity,omitempty,indexed"`
		LowestAsk  float64    `json:"lowest_ask"  db:"lowest_ask,omitempty,indexed"`
		HighestBid float64    `json:"highest_bid" db:"highest_bid,omitempty,indexed"`
		RecentAsk  *time.Time `json:"recent_ask"  db:"recent_ask,omitempty,indexed"`
		UpdatedAt  *time.Time `json:"updated_at"  db:"updated_at,omitempty,indexed"`
	}

	// CatalogService provides access to catalog service.
	CatalogService interface {
		// Catalogs returns a list of catalogs.
		Catalogs(opts FindOpts) ([]Catalog, *FindMetadata, error)

		// Catalog returns catalog details by item ID.
		Catalog(itemID string) (*Catalog, error)

		// Index creates or update index entry using item ID.
		Index(itemID string) (*Catalog, error)
	}

	CatalogStorage interface {
		// Find returns a list of catalogs from data store.
		Find(opts FindOpts) ([]Catalog, error)

		// Count returns number of catalog from data store.
		Count(FindOpts) (int, error)

		// Get returns catalog details by id from data store.
		Get(id string) (*Catalog, error)

		// Index persists a new catalog to data store.
		Index(itemID string) (*Catalog, error)
	}
)
