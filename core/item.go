package core

import (
	"context"
	"strconv"
	"time"
)

// User error types.
const (
	ItemErrNotFound Errors = iota + 2000
	ItemErrRequiredID
	ItemErrRequiredFields
	ItemErrProfileInvalidStatus
	ItemErrProfileNotesLimit
)

// sets error text definition.
func init() {
	appErrorText[ItemErrNotFound] = "item not found"
	appErrorText[ItemErrRequiredID] = "item id is required"
	appErrorText[ItemErrRequiredFields] = "item fields are required"
	appErrorText[ItemErrProfileInvalidStatus] = "item status not allowed"
	appErrorText[ItemErrProfileNotesLimit] = "item notes text limit reached"
}

const maxItemNotesLen = 100

// Item statuses.
const (
	ItemStatusPending  ItemStatus = 100
	ItemStatusLive     ItemStatus = 200
	ItemStatusReserved ItemStatus = 300
	ItemStatusSold     ItemStatus = 400
	ItemStatusRemoved  ItemStatus = 500
)

type (
	// ItemStatus represents item status.
	ItemStatus uint

	// Item represents item information.
	Item struct {
		ID        string     `json:"id"         db:"id,omitempty"`
		Name      string     `json:"item"       db:"item,omitempty"        valid:"required"`
		Hero      string     `json:"hero"       db:"hero,omitempty"        valid:"required"`
		Price     float64    `json:"price"      db:"price,omitempty"       valid:"required"`
		Currency  string     `json:"currency"   db:"currency,omitempty"`
		Notes     string     `json:"notes"      db:"notes,omitempty"`
		Status    ItemStatus `json:"status"     db:"status,omitempty"`
		CreatedAt *time.Time `json:"created_at" db:"created_at,omitempty"`
		UpdatedAt *time.Time `json:"updated_at" db:"updated_at,omitempty"`
		// Include related fields.
		User *User `json:"user,omitempty" db:"-"`
	}

	// ItemService provides access to item service.
	ItemService interface {
		// Items returns a list of items.
		Items(opts FindOpts) ([]Item, error)

		// Item returns item details by id.
		Item(id string) (*Item, error)

		// Create saves new item details.
		Create(context.Context, *Item) error

		// Update saves item details changes.
		Update(context.Context, *Item) error
	}

	ItemStorage interface {
		// Find returns a list of items from data store.
		Find(opts FindOpts) ([]Item, error)

		// Get returns item details by id from data store.
		Get(id string) (*Item, error)

		// Create persists a new item to data store.
		Create(*Item) error

		// Update persists item changes to data store.
		Update(*Item) error
	}
)

var itemStatusTexts = map[ItemStatus]string{
	ItemStatusPending:  "pending",
	ItemStatusLive:     "live",
	ItemStatusReserved: "reserved",
	ItemStatusSold:     "sold",
	ItemStatusRemoved:  "removed",
}

// CheckCreate validates field on creating new item.
func (i Item) CheckCreate() error {
	// Check required fields.
	if err := validator.Struct(i); err != nil {
		return err
	}

	// Check title and description length.
	if len(i.Notes) > maxItemNotesLen {
		return ItemErrProfileNotesLimit
	}

	return nil
}

// SetDefault sets default values for a new item.
func (i Item) SetDefaults() Item {
	i.Status = ItemStatusLive
	i.Currency = "USD"
	return i
}

// String returns text value of a post status.
func (s ItemStatus) String() string {
	t, ok := itemStatusTexts[s]
	if !ok {
		return strconv.Itoa(int(s))
	}

	return t
}
