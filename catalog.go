package dgx

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
		ID           string   `json:"id"         db:"id,omitempty"`
		Slug         string   `json:"slug"       db:"slug,omitempty,indexed"`
		Name         string   `json:"name"       db:"name,omitempty,indexed"`
		Hero         string   `json:"hero"       db:"hero,omitempty,indexed"`
		Image        string   `json:"image"      db:"image,omitempty"`
		Origin       string   `json:"origin"     db:"origin,omitempty,indexed"`
		Rarity       string   `json:"rarity"     db:"rarity,omitempty,indexed"`
		Contributors []string `json:"-"          db:"contributors,omitempty"`
		ViewCount    int      `json:"view_count" db:"view_count,omitempty,indexed"`
		// Market summary details.
		Quantity      int        `json:"quantity"       db:"quantity,omitempty"`
		LowestAsk     float64    `json:"lowest_ask"     db:"lowest_ask,omitempty"`
		MedianAsk     float64    `json:"median_ask"     db:"median_ask,omitempty"`
		RecentAsk     *time.Time `json:"recent_ask"     db:"recent_ask,omitempty,indexed"`
		HighestBid    float64    `json:"highest_bid"    db:"highest_bid,omitempty"`
		RecentBid     *time.Time `json:"recent_bid"     db:"recent_bid,omitempty,indexed"`
		BidCount      int        `json:"bid_count"      db:"bid_count,omitempty"`
		ReservedCount int        `json:"reserved_count" db:"reserved_count,omitempty"`
		SoldCount     int        `json:"sold_count"     db:"sold_count,omitempty"`
		// Sale summary details are derived from reserved and sold status and not the same as sold.
		SaleCount  int        `json:"sale_count"  db:"sale_count,omitempty"`
		AvgSale    float64    `json:"avg_sale"    db:"avg_sale,omitempty"`
		RecentSale *time.Time `json:"recent_sale" db:"recent_sale,omitempty"`
		CreatedAt  *time.Time `json:"created_at"  db:"created_at,omitempty,indexed"`
		UpdatedAt  *time.Time `json:"updated_at"  db:"updated_at,omitempty,indexed"`
		// Include related fields.
		Asks []Market `json:"asks" db:"-"`
		Bids []Market `json:"bids" db:"-"`
	}

	// CatalogStorage defines operation for market indexed items.
	CatalogStorage interface {
		// Find returns a list of catalogs from data store.
		Find(opts FindOpts) ([]Catalog, error)

		// Count returns number of catalog from data store.
		Count(FindOpts) (int, error)

		// Get returns catalog details by id from data store.
		Get(id string) (*Catalog, error)

		// Index persists a new catalog to data store.
		Index(itemID string) (*Catalog, error)

		// Trending returns a list if top 10 trending catalog.
		Trending() ([]Catalog, error)
	}
)
