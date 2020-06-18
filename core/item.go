package core

import (
	"context"
	"time"
)

// Item error types.
const (
	ItemErrNotFound Errors = iota + 2000
	ItemErrRequiredID
	ItemErrRequiredFields
)

// sets error text definition.
func init() {
	appErrorText[ItemErrNotFound] = "item not found"
	appErrorText[ItemErrRequiredID] = "item id is required"
	appErrorText[ItemErrRequiredFields] = "item fields are required"
}

type (
	// ItemStatus represents item status.
	ItemStatus uint

	// Item represents item information.
	Item struct {
		ID        string     `json:"id"         db:"id,omitempty"`
		Name      string     `json:"item"       db:"item,omitempty"        valid:"required"`
		Hero      string     `json:"hero"       db:"hero,omitempty"        valid:"required"`
		Tag       string     `json:"tag"        db:"tag,omitempty"`
		CreatedAt *time.Time `json:"created_at" db:"created_at,omitempty"`
		UpdatedAt *time.Time `json:"updated_at" db:"updated_at,omitempty"`
	}

	// ItemService provides access to item service.
	ItemService interface {
		// Items returns a list of items.
		Items(opts FindOpts) ([]Item, FindMetadata, error)

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

// CheckCreate validates field on creating new item.
func (i Item) CheckCreate() error {
	// Check required fields.
	if err := validator.Struct(i); err != nil {
		return err
	}

	return nil
}

// SetDefault sets default values for a new item.
func (i Item) SetDefaults() Item {
	return i
}
