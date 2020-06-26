package core

import (
	"context"
	"math"
	"strconv"
	"time"
)

// Market error types.
const (
	MarketErrNotFound Errors = iota + 2100
	MarketErrRequiredID
	MarketErrRequiredFields
	MarketErrInvalidStatus
	MarketErrNotesLimit
)

// sets error text definition.
func init() {
	appErrorText[MarketErrNotFound] = "market not found"
	appErrorText[MarketErrRequiredID] = "market id is required"
	appErrorText[MarketErrRequiredFields] = "market fields are required"
	appErrorText[MarketErrInvalidStatus] = "market status not allowed"
	appErrorText[MarketErrNotesLimit] = "market notes text limit reached"
}

const maxMarketNotesLen = 120

// Market statuses.
const (
	MarketStatusPending  MarketStatus = 100
	MarketStatusLive     MarketStatus = 200
	MarketStatusReserved MarketStatus = 300
	MarketStatusSold     MarketStatus = 400
	MarketStatusRemoved  MarketStatus = 500
)

type (
	// MarketStatus represents market status.
	MarketStatus uint

	// Market represents market information.
	Market struct {
		ID        string       `json:"id"         db:"id,omitempty"`
		UserID    string       `json:"user_id"    db:"user_id,omitempty"     valid:"required"`
		ItemID    string       `json:"item_id"    db:"item_id,omitempty"     valid:"required"`
		Price     float64      `json:"price"      db:"price,omitempty"       valid:"required"`
		Currency  string       `json:"currency"   db:"currency,omitempty"`
		Notes     string       `json:"notes"      db:"notes,omitempty"`
		Status    MarketStatus `json:"status"     db:"status,omitempty"`
		CreatedAt *time.Time   `json:"created_at" db:"created_at,omitempty"`
		UpdatedAt *time.Time   `json:"updated_at" db:"updated_at,omitempty"`
		// Include related fields.
		User *User `json:"user,omitempty" db:"-"`
		Item *Item `json:"item,omitempty" db:"-"`
	}

	// MarketIndex represents aggregation of market entries.
	MarketIndex struct {
		ItemID     string     `json:"item_id"     db:"item_id,omitempty"`
		Quantity   int        `json:"quantity"    db:"quantity,omitempty"`
		LowestAsk  float64    `json:"lowest_ask"  db:"lowest_ask,omitempty"`
		HighestBid float64    `json:"highest_bid" db:"highest_bid,omitempty"`
		RecentAsk  *time.Time `json:"recent_ask"  db:"recent_ask,omitempty"`
		// Include related fields.
		Item `json:"item,omitempty"`
	}

	// MarketService provides access to market service.
	MarketService interface {
		// Markets returns a list of markets.
		Markets(ctx context.Context, opts FindOpts) ([]Market, *FindMetadata, error)

		// Market returns market details by id.
		Market(ctx context.Context, id string) (*Market, error)

		// Create saves new market details.
		Create(context.Context, *Market) error

		// Update saves market details changes.
		Update(context.Context, *Market) error

		// Index returns a list of indexed markets.
		Index(opts FindOpts) ([]MarketIndex, *FindMetadata, error)
	}

	MarketStorage interface {
		// Find returns a list of markets from data store.
		Find(opts FindOpts) ([]Market, error)

		// Count returns number of market from data store.
		Count(FindOpts) (int, error)

		// Get returns market details by id from data store.
		Get(id string) (*Market, error)

		// Create persists a new market to data store.
		Create(*Market) error

		// Update persists market changes to data store.
		Update(*Market) error

		// Find returns a list o aggregated market index from data store.
		FindIndex(opts FindOpts) ([]MarketIndex, error)

		// Count returns number of aggregated market index from data store.
		CountIndex(FindOpts) (int, error)
	}
)

var MarketStatusTexts = map[MarketStatus]string{
	MarketStatusPending:  "pending",
	MarketStatusLive:     "live",
	MarketStatusReserved: "reserved",
	MarketStatusSold:     "sold",
	MarketStatusRemoved:  "removed",
}

// CheckCreate validates field on creating new market.
func (i Market) CheckCreate() error {
	// Check required fields.
	if err := validator.Struct(i); err != nil {
		return err
	}

	// Check market notes length.
	if len(i.Notes) > maxMarketNotesLen {
		return MarketErrNotesLimit
	}

	return nil
}

// CheckUpdate validates field on updating market.
func (i Market) CheckUpdate() error {
	if i.Notes != "" && len(i.Notes) > maxMarketNotesLen {
		return MarketErrNotesLimit
	}

	_, ok := MarketStatusTexts[i.Status]
	if i.Status != 0 && !ok {
		return MarketErrInvalidStatus
	}

	return nil
}

const defaultCurrency = "USD"

// SetDefault sets default values for a new market.
func (i *Market) SetDefaults() {
	i.Status = MarketStatusLive
	i.Currency = defaultCurrency
	i.Price = priceToTenths(i.Price)
}

// String returns text value of a post status.
func (s MarketStatus) String() string {
	t, ok := MarketStatusTexts[s]
	if !ok {
		return strconv.Itoa(int(s))
	}

	return t
}

func priceToTenths(n float64) float64 {
	const dec = 100
	return math.Round(n*dec) / dec
}
