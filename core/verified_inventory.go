package core

import (
	"context"
	"time"
)

// Inventory error types.
const (
	InventoryErrNotFound Errors = iota + 6100
	InventoryErrRequiredID
	InventoryErrRequiredFields
)

// sets error text definition.
func init() {
	appErrorText[InventoryErrNotFound] = "report not found"
	appErrorText[InventoryErrRequiredID] = "report id is required"
	appErrorText[InventoryErrRequiredFields] = "report fields are required"
}

// Inventory statuses.
const (
	// InventoryStatusNoHit buyer's inventory successfully parsed
	// but the item did not find any in match.
	InventoryStatusNoHit InventoryStatus = 100

	// InventoryStatusNameVerified item exists on buyer's inventory
	// base on the item name challenge.
	//
	// No-gift info might mean:
	// 1. Buyer cleared the gift information
	// 2. Buyer is the original owner of the item
	// 3. Item might come from another source
	InventoryStatusNameVerified InventoryStatus = 200

	// InventoryStatusSenderVerified both item existence and gift
	// information matched the seller's avatar name. We could
	// also use the date received to check against delivery data
	// to strengthen its validity.
	InventoryStatusSenderVerified InventoryStatus = 300

	// InventoryStatusPrivate buyer's inventory is not visible to
	// public and we can do nothing about it.
	InventoryStatusPrivate InventoryStatus = 400

	// InventoryStatusError error occurred during API request or
	// parsing inventory error.
	InventoryStatusError InventoryStatus = 500
)

type (

	// InventoryStatus represents delivery status.
	InventoryStatus uint

	/// Inventory represents steam inventory delivery.
	Inventory struct {
		ID               string          `json:"id"                 db:"id,omitempty,omitempty"`
		MarketID         string          `json:"market_id"          db:"market_id,omitempty"`
		BuyerConfirmed   *bool           `json:"buyer_confirmed"    db:"buyer_confirmed,omitempty"`
		BuyerConfirmedAt *time.Time      `json:"buyer_confirmed_at" db:"buyer_confirmed_at,omitempty"`
		Status           InventoryStatus `json:"status"             db:"status,omitempty"`
		Assets           []SteamAsset    `json:"steam_assets"       db:"steam_assets,omitempty"`
		Retries          int             `json:"retries"            db:"retries,omitempty"`
		CreatedAt        *time.Time      `json:"created_at"         db:"created_at,omitempty,indexed,omitempty"`
		UpdatedAt        *time.Time      `json:"updated_at"         db:"updated_at,omitempty,indexed,omitempty"`
	}

	// InventoryService provides access to Inventory service.
	InventoryService interface {
		// Deliveries returns a list of inventories.
		Deliveries(opts FindOpts) ([]Inventory, *FindMetadata, error)

		// Inventory returns Inventory details by id.
		Inventory(id string) (*Inventory, error)

		// Create saves new Inventory details.
		Set(context.Context, *Inventory) error
	}

	// InventoryStorage defines operation for Inventory records.
	InventoryStorage interface {
		// Find returns a list of inventories from data store.
		Find(opts FindOpts) ([]Inventory, error)

		// Count returns number of inventories from data store.
		Count(FindOpts) (int, error)

		// Get returns Inventory details by id from data store.
		Get(id string) (*Inventory, error)

		// Create persists a new Inventory to data store.
		Create(*Inventory) error

		// Update save changes of Inventory to data store.
		Update(*Inventory) error
	}
)
