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
	MarketErrInvalidPrice
	MarketErrQtyLimitPerUser
	MarketErrRequiredPartnerURL
)

// sets error text definition.
func init() {
	appErrorText[MarketErrNotFound] = "market not found"
	appErrorText[MarketErrRequiredID] = "market id is required"
	appErrorText[MarketErrRequiredFields] = "market fields are required"
	appErrorText[MarketErrInvalidStatus] = "market status not allowed"
	appErrorText[MarketErrNotesLimit] = "market notes text limit reached"
	appErrorText[MarketErrInvalidPrice] = "market price is invalid"
	appErrorText[MarketErrQtyLimitPerUser] = "market quantity limit(5) per item reached"
	appErrorText[MarketErrRequiredPartnerURL] = "market partner steam url is required"
}

const (
	maxMarketNotesLen        = 200
	MaxMarketQtyLimitPerUser = 5

	// Trend scoring rates use for trend ranking.
	TrendScoreRateView        = 0.05
	TrendScoreRateMarketEntry = 0.01
	TrendScoreRateReserved    = 4
	TrendScoreRateSold        = 4
)

// Market statuses.
const (
	MarketStatusPending   MarketStatus = 100
	MarketStatusLive      MarketStatus = 200
	MarketStatusReserved  MarketStatus = 300
	MarketStatusSold      MarketStatus = 400
	MarketStatusRemoved   MarketStatus = 500
	MarketStatusCancelled MarketStatus = 600
)

type (
	// MarketStatus represents market status.
	MarketStatus uint

	// Market represents market information.
	Market struct {
		ID             string       `json:"id"               db:"id,omitempty"`
		UserID         string       `json:"user_id"          db:"user_id,omitempty,indexed"   valid:"required"`
		ItemID         string       `json:"item_id"          db:"item_id,omitempty,indexed"   valid:"required"`
		Status         MarketStatus `json:"status"           db:"status,omitempty,indexed"    valid:"required"`
		Price          float64      `json:"price"            db:"price,omitempty,indexed"     valid:"required"`
		Currency       string       `json:"currency"         db:"currency,omitempty"`
		PartnerSteamID string       `json:"partner_steam_id" db:"partner_steam_id,omitempty"`
		Notes          string       `json:"notes"            db:"notes,omitempty"`
		CreatedAt      *time.Time   `json:"created_at"       db:"created_at,omitempty,indexed"`
		UpdatedAt      *time.Time   `json:"updated_at"       db:"updated_at,omitempty,indexed"`
		// Include related fields.
		User *User `json:"user,omitempty" db:"user,omitempty"`
		Item *Item `json:"item,omitempty" db:"item,omitempty"`
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

		// Catalog returns a list of catalogs.
		Catalog(opts FindOpts) ([]Catalog, *FindMetadata, error)

		// CatalogDetails returns catalog details by item id.
		CatalogDetails(id string) (*Catalog, error)

		// TrendingCatalog returns a top 10 trending catalogs.
		TrendingCatalog(opts FindOpts) ([]Catalog, *FindMetadata, error)
	}

	// MarketStorage defines operation for market records.
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
	}
)

var MarketStatusTexts = map[MarketStatus]string{
	MarketStatusPending:   "pending",
	MarketStatusLive:      "live",
	MarketStatusReserved:  "reserved",
	MarketStatusSold:      "sold",
	MarketStatusRemoved:   "removed",
	MarketStatusCancelled: "cancelled",
}

// CheckCreate validates field on creating new market.
func (m Market) CheckCreate() error {
	// Check required fields.
	if err := validator.Struct(m); err != nil {
		return err
	}

	// Check valid market price.
	if m.Price <= 0 {
		return MarketErrInvalidPrice
	}

	// Check market notes length.
	if len(m.Notes) > maxMarketNotesLen {
		return MarketErrNotesLimit
	}

	return nil
}

// CheckUpdate validates field on updating market.
func (m Market) CheckUpdate() error {
	if m.Notes != "" && len(m.Notes) > maxMarketNotesLen {
		return MarketErrNotesLimit
	}

	_, ok := MarketStatusTexts[m.Status]
	if m.Status != 0 && !ok {
		return MarketErrInvalidStatus
	}

	if m.Status == MarketStatusReserved && m.PartnerSteamID == "" {
		return MarketErrRequiredPartnerURL
	}

	return nil
}

const defaultCurrency = "USD"

// SetDefault sets default values for a new market.
func (m *Market) SetDefaults() {
	m.Status = MarketStatusLive
	m.Currency = defaultCurrency
	m.Price = priceToTenths(m.Price)
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
