package dgx

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
	MarketErrInvalidBidPrice
	MarketErrInvalidAskPrice
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
	appErrorText[MarketErrInvalidBidPrice] = "market bid should be lower than lowest ask price"
	appErrorText[MarketErrInvalidAskPrice] = "market ask should be higher than highest bid price"
}

const (
	maxMarketNotesLen               = 200
	MaxMarketQtyLimitPerFreeUser    = 1
	MaxMarketQtyLimitPerPremiumUser = 5

	MarketAskExpirationDays = 30
	MarketBidExpirationDays = 7

	MarketSweepExpiredDays = 30
	MarketSweepRemovedDays = 60
)

// Market types.
const (
	MarketTypeAsk MarketType = 10 // default
	MarketTypeBid MarketType = 20
)

// Market statuses.
const (
	MarketStatusPending      MarketStatus = 100
	MarketStatusLive         MarketStatus = 200
	MarketStatusReserved     MarketStatus = 300
	MarketStatusSold         MarketStatus = 400
	MarketStatusBidCompleted MarketStatus = 410
	MarketStatusRemoved      MarketStatus = 500
	MarketStatusCancelled    MarketStatus = 600
	MarketStatusExpired      MarketStatus = 700
)

// Market trending score rates.
const (
	TrendScoreRateView        = 0.05
	TrendScoreRateMarketEntry = 0.01
	TrendScoreRateReserved    = 4
	TrendScoreRateSold        = 4
	TrendScoreRateBid         = 2
)

type (
	// MarketType represents market type.
	MarketType uint

	// MarketStatus represents market status.
	MarketStatus uint

	// Market represents market information.
	Market struct {
		ID             string       `json:"id"               db:"id,omitempty"`
		UserID         string       `json:"user_id"          db:"user_id,omitempty,indexed"   valid:"required"`
		ItemID         string       `json:"item_id"          db:"item_id,omitempty,indexed"   valid:"required"`
		Type           MarketType   `json:"type"             db:"type,omitempty,indexed"      valid:"required"`
		Status         MarketStatus `json:"status"           db:"status,omitempty,indexed"    valid:"required"`
		Price          float64      `json:"price"            db:"price,omitempty,indexed"     valid:"required"`
		Currency       string       `json:"currency"         db:"currency,omitempty"`
		PartnerSteamID string       `json:"partner_steam_id" db:"partner_steam_id,omitempty"`
		Notes          string       `json:"notes"            db:"notes,omitempty"`
		CreatedAt      *time.Time   `json:"created_at"       db:"created_at,omitempty,indexed"`
		UpdatedAt      *time.Time   `json:"updated_at"       db:"updated_at,omitempty,indexed"`

		// Will be use for full-text searching.
		SearchText    string `json:"-"               db:"search_text,omitempty,indexed"`
		UserRankScore int    `json:"user_rank_score" db:"user_rank_score,omitempty,indexed"`

		InventoryStatus InventoryStatus `json:"inventory_status" db:"inventory_status,omitempty,indexed"`
		DeliveryStatus  DeliveryStatus  `json:"delivery_status"  db:"delivery_status,omitempty,indexed"`

		// Include related fields.
		User      *User      `json:"user,omitempty"      db:"user,omitempty"`
		Item      *Item      `json:"item,omitempty"      db:"item,omitempty"`
		Delivery  *Delivery  `json:"delivery,omitempty"  db:"delivery,omitempty"`
		Inventory *Inventory `json:"inventory,omitempty" db:"inventory,omitempty"`

		// NOTE! Experimental for reselling feature.
		Resell        *bool  `json:"resell"          db:"resell,omitempty"`
		SellerSteamID string `json:"seller_steam_id" db:"seller_steam_id,omitempty"`
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

		// UpdateUserRankScore sets new user ranking score on all live market by user id.
		UpdateUserRankScore(userID string) error

		// Index composes market data for faster search and retrieval.
		//Index(ctx context.Context, id string) (*Market, error)

		// AutoCompleteBid detects if there's matching reservation on buy order and automatically
		// resolve it by setting complete-bid status.
		AutoCompleteBid(ctx context.Context, ask Market, partnerSteamID string) error

		// Catalog returns a list of catalogs.
		Catalog(opts FindOpts) ([]Catalog, *FindMetadata, error)

		// CatalogDetails returns catalog details by item id.
		CatalogDetails(id string, opts FindOpts) (*Catalog, error)

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

		// BaseUpdate persists market changes to data store and
		// will not update updated_at field.
		BaseUpdate(*Market) error

		// PendingInventoryStatus returns market entries that is pending for checking
		// inventory status or needs re-processing of re-process error status.
		PendingInventoryStatus(o FindOpts) ([]Market, error)

		// PendingDeliveryStatus returns market entries that is pending for checking
		// delivery status or needs re-processing of re-process error status.
		PendingDeliveryStatus(o FindOpts) ([]Market, error)

		RevalidateDeliveryStatus(o FindOpts) ([]Market, error)

		// Index composes market data for faster search and retrieval.
		Index(id string) (*Market, error)

		// UpdateUserScore sets new rank score value of all live market by user ID.
		UpdateUserScore(userID string, rankScore int) error

		// UpdateExpiring sets live items to expired status by expiration time.
		UpdateExpiring(t MarketType, b UserBoon, expiration time.Time) (itemIDs []string, err error)

		BulkDeleteByStatus(ms MarketStatus, cutOff time.Time, limit int) error

		UpdateExpiringResell(b UserBoon) (itemIDs []string, err error)
	}
)

var MarketStatusTexts = map[MarketStatus]string{
	MarketStatusPending:      "pending",
	MarketStatusLive:         "live",
	MarketStatusReserved:     "reserved",
	MarketStatusSold:         "sold",
	MarketStatusBidCompleted: "completed",
	MarketStatusRemoved:      "removed",
	MarketStatusCancelled:    "cancelled",
	MarketStatusExpired:      "expired",
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

// SetDefaults sets default values for a new market.
func (m Market) SetDefaults() *Market {
	m.Status = MarketStatusLive
	m.Currency = defaultCurrency
	m.Price = priceToTenths(m.Price)
	if m.Type == 0 {
		m.Type = MarketTypeAsk
	}
	return &m
}

// IsResell check if the market is re-sell item.
func (m Market) IsResell() bool {
	return m.Type == MarketTypeAsk && m.Resell != nil && *m.Resell
}

// String returns text value of a market status.
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
