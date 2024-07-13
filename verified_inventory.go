package dgx

import (
	"context"
	"strconv"
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

	// InventoryStatusVerified item exists on inventory base on
	// the item name challenge.
	InventoryStatusVerified InventoryStatus = 200

	// InventoryStatusPrivate buyer's inventory is not visible to
	// public, and we can do nothing about it.
	InventoryStatusPrivate InventoryStatus = 400

	// InventoryStatusError error occurred during API request or
	// parsing inventory error.
	InventoryStatusError InventoryStatus = 500
)

var inventoryStatusTexts = map[InventoryStatus]string{
	InventoryStatusNoHit:    "no hit",
	InventoryStatusVerified: "verified",
	InventoryStatusPrivate:  "private",
	InventoryStatusError:    "error",
}

type (

	// InventoryStatus represents inventory status.
	InventoryStatus uint

	/// Inventory represents steam inventory inventory.
	Inventory struct {
		ID          string          `json:"id"                 db:"id,omitempty,omitempty"`
		MarketID    string          `json:"market_id"          db:"market_id,omitempty,indexed" valid:"required"`
		Status      InventoryStatus `json:"status"             db:"status,omitempty,indexed"    valid:"required"`
		Assets      []SteamAsset    `json:"steam_assets"       db:"steam_assets,omitempty"`
		Retries     int             `json:"retries"            db:"retries,omitempty"`
		BundleCount int             `json:"bundle_count"       db:"bundle_count,omitempty"`
		CreatedAt   *time.Time      `json:"created_at"         db:"created_at,omitempty,indexed,omitempty"`
		UpdatedAt   *time.Time      `json:"updated_at"         db:"updated_at,omitempty,indexed,omitempty"`
	}

	// InventoryService provides access to Inventory service.
	InventoryService interface {
		// Inventories returns a list of deliveries.
		Inventories(opts FindOpts) ([]Inventory, *FindMetadata, error)

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

		// GetByMarketID returns Inventory details by market id from data store.
		GetByMarketID(marketID string) (*Inventory, error)

		// Create persists a new Inventory to data store.
		Create(*Inventory) error

		// Update save changes of Inventory to data store.
		Update(*Inventory) error
	}
)

// CheckCreate validates field on creating new inventory.
func (i Inventory) CheckCreate() error {
	// Check required fields.
	if err := validator.Struct(i); err != nil {
		return err
	}

	return nil
}

func (i Inventory) CountBundles() (total int) {
	if i.Assets == nil {
		return
	}

	for _, aa := range i.Assets {
		if aa.IsBundled() {
			total++
		}
	}
	return total
}

// String returns text value of a inventory status.
func (s InventoryStatus) String() string {
	t, ok := inventoryStatusTexts[s]
	if !ok {
		return strconv.Itoa(int(s))
	}

	return t
}

// RetriesExceeded when it reached 5 reties.
func (d Inventory) RetriesExceeded() bool {
	return d.Retries > 3
}
