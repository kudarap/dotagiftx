package core

import (
	"context"
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
	// Catalog represents indexed market.
	Catalog struct {
		ItemID     string     `json:"item_id"     db:"item_id,omitempty"`
		Quantity   int        `json:"quantity"    db:"quantity,omitempty"`
		LowestAsk  float64    `json:"lowest_ask"  db:"lowest_ask,omitempty"`
		HighestBid float64    `json:"highest_bid" db:"highest_bid,omitempty"`
		RecentAsk  *time.Time `json:"recent_ask"  db:"recent_ask,omitempty"`
		// Include related fields.
		Item
	}

	// CatalogService provides access to catalog service.
	CatalogService interface {
		// Catalogs returns a list of catalogs.
		Catalogs(ctx context.Context, opts FindOpts) ([]Catalog, *FindMetadata, error)

		// Catalog returns catalog details by item ID.
		Catalog(itemID string) (*Catalog, error)

		// Index creates or update index entry using item ID.
		Index(itemID string) ([]Catalog, *FindMetadata, error)
	}

	CatalogStorage interface {
		// Find returns a list of catalogs from data store.
		Find(opts FindOpts) ([]Catalog, error)

		// Count returns number of catalog from data store.
		Count(FindOpts) (int, error)

		// Get returns catalog details by id from data store.
		Get(id string) (*Catalog, error)

		// Index persists a new catalog to data store.
		Index(*Catalog) error
	}
)
