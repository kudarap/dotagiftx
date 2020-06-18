package core

import (
	"context"
	"strconv"
	"time"
)

// Sell error types.
const (
	SellErrNotFound Errors = iota + 2100
	SellErrRequiredID
	SellErrRequiredFields
	SellErrInvalidStatus
	SellErrNotesLimit
)

// sets error text definition.
func init() {
	appErrorText[SellErrNotFound] = "sell not found"
	appErrorText[SellErrRequiredID] = "sell id is required"
	appErrorText[SellErrRequiredFields] = "sell fields are required"
	appErrorText[SellErrInvalidStatus] = "sell status not allowed"
	appErrorText[SellErrNotesLimit] = "sell notes text limit reached"
}

const maxSellNotesLen = 120

// Sell statuses.
const (
	SellStatusPending  SellStatus = 100
	SellStatusLive     SellStatus = 200
	SellStatusReserved SellStatus = 300
	SellStatusSold     SellStatus = 400
	SellStatusRemoved  SellStatus = 500
)

type (
	// SellStatus represents sell status.
	SellStatus uint

	// Sell represents sell information.
	Sell struct {
		ID        string     `json:"id"         db:"id,omitempty"`
		UserID    string     `json:"user_id"    db:"user_id,omitempty"     valid:"required"`
		ItemID    string     `json:"item_id"    db:"item_id,omitempty"     valid:"required"`
		Price     float64    `json:"price"      db:"price,omitempty"       valid:"required"`
		Currency  string     `json:"currency"   db:"currency,omitempty"`
		Notes     string     `json:"notes"      db:"notes,omitempty"`
		Status    SellStatus `json:"status"     db:"status,omitempty"`
		CreatedAt *time.Time `json:"created_at" db:"created_at,omitempty"`
		UpdatedAt *time.Time `json:"updated_at" db:"updated_at,omitempty"`
		// Include related fields.
		User *User `json:"user,omitempty" db:"-"`
		Item *Item `json:"item,omitempty" db:"-"`
	}

	// SellService provides access to sell service.
	SellService interface {
		// Sells returns a list of sells.
		Sells(opts FindOpts) ([]Sell, *FindMetadata, error)

		// Sell returns sell details by id.
		Sell(id string) (*Sell, error)

		// Create saves new sell details.
		Create(context.Context, *Sell) error

		// Update saves sell details changes.
		Update(context.Context, *Sell) error
	}

	SellStorage interface {
		// Find returns a list of sells from data store.
		Find(opts FindOpts) ([]Sell, error)

		// Count returns number of sell from data store.
		Count(FindOpts) (int, error)

		// Get returns sell details by id from data store.
		Get(id string) (*Sell, error)

		// Create persists a new sell to data store.
		Create(*Sell) error

		// Update persists sell changes to data store.
		Update(*Sell) error
	}
)

var sellStatusTexts = map[SellStatus]string{
	SellStatusPending:  "pending",
	SellStatusLive:     "live",
	SellStatusReserved: "reserved",
	SellStatusSold:     "sold",
	SellStatusRemoved:  "removed",
}

// CheckCreate validates field on creating new sell.
func (i Sell) CheckCreate() error {
	// Check required fields.
	if err := validator.Struct(i); err != nil {
		return err
	}

	// Check sell notes length.
	if len(i.Notes) > maxSellNotesLen {
		return SellErrNotesLimit
	}

	return nil
}

// CheckUpdate validates field on updating sell.
func (i Sell) CheckUpdate() error {
	if i.Notes != "" && len(i.Notes) > maxSellNotesLen {
		return SellErrNotesLimit
	}

	_, ok := sellStatusTexts[i.Status]
	if i.Status != 0 && !ok {
		return SellErrInvalidStatus
	}

	return nil
}

const defaultCurrency = "USD"

// SetDefault sets default values for a new sell.
func (i *Sell) SetDefaults() {
	i.Status = SellStatusLive
	i.Currency = defaultCurrency
}

// String returns text value of a post status.
func (s SellStatus) String() string {
	t, ok := sellStatusTexts[s]
	if !ok {
		return strconv.Itoa(int(s))
	}

	return t
}
