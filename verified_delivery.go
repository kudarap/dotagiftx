package dgx

import (
	"context"
	"strconv"
	"time"
)

// Delivery error types.
const (
	DeliveryErrNotFound Errors = iota + 6000
	DeliveryErrRequiredID
	DeliveryErrRequiredFields
)

// sets error text definition.
func init() {
	appErrorText[DeliveryErrNotFound] = "report not found"
	appErrorText[DeliveryErrRequiredID] = "report id is required"
	appErrorText[DeliveryErrRequiredFields] = "report fields are required"
}

// DeliveryRetryLimit max retry to process verification.
const DeliveryRetryLimit = 30

// Delivery statuses.
const (
	// DeliveryStatusNoHit buyer's inventory successfully parsed
	// but the item did not find any in match.
	DeliveryStatusNoHit DeliveryStatus = 100

	// DeliveryStatusNameVerified item exists on buyer's inventory
	// base on the item name challenge.
	//
	// No-gift info might mean:
	// 1. Buyer cleared the gift information
	// 2. Buyer is the original owner of the item
	// 3. Item might come from another source
	DeliveryStatusNameVerified DeliveryStatus = 200

	// DeliveryStatusSenderVerified both item existence and gift
	// information matched the seller's avatar name. We could
	// also use the date received to check against delivery data
	// to strengthen its validity.
	DeliveryStatusSenderVerified DeliveryStatus = 300

	// DeliveryStatusPrivate buyer's inventory is not visible to
	// public, and we can do nothing about it.
	DeliveryStatusPrivate DeliveryStatus = 400

	// DeliveryStatusError error occurred during API request or
	// parsing inventory error.
	DeliveryStatusError DeliveryStatus = 500
)

var deliveryStatusTexts = map[DeliveryStatus]string{
	DeliveryStatusNoHit:          "no hit",
	DeliveryStatusNameVerified:   "name verified",
	DeliveryStatusSenderVerified: "sender verified",
	DeliveryStatusPrivate:        "private",
	DeliveryStatusError:          "error",
}

type (

	// DeliveryStatus represents delivery status.
	DeliveryStatus uint

	// Delivery represents steam inventory delivery.
	Delivery struct {
		ID               string         `json:"id"                 db:"id,omitempty,omitempty"`
		MarketID         string         `json:"market_id"          db:"market_id,omitempty,indexed" valid:"required"`
		BuyerConfirmed   *bool          `json:"buyer_confirmed"    db:"buyer_confirmed,omitempty"`
		BuyerConfirmedAt *time.Time     `json:"buyer_confirmed_at" db:"buyer_confirmed_at,omitempty"`
		GiftOpened       *bool          `json:"gift_opened"        db:"gift_opened,omitempty"`
		Status           DeliveryStatus `json:"status"             db:"status,omitempty,indexed"    valid:"required"`
		Assets           []SteamAsset   `json:"steam_assets"       db:"steam_assets,omitempty"`
		Retries          int            `json:"retries"            db:"retries,omitempty"`
		CreatedAt        *time.Time     `json:"created_at"         db:"created_at,omitempty,indexed,omitempty"`
		UpdatedAt        *time.Time     `json:"updated_at"         db:"updated_at,omitempty,indexed,omitempty"`
	}

	// DeliveryService provides access to Delivery service.
	DeliveryService interface {
		// Deliveries returns a list of deliveries.
		Deliveries(opts FindOpts) ([]Delivery, *FindMetadata, error)

		// Delivery returns Delivery details by id.
		Delivery(id string) (*Delivery, error)

		// Set saves new Delivery details.
		Set(context.Context, *Delivery) error
	}

	// DeliveryStorage defines operation for Delivery records.
	DeliveryStorage interface {
		// Find returns a list of deliveries from data store.
		Find(opts FindOpts) ([]Delivery, error)

		// Count returns number of deliveries from data store.
		Count(FindOpts) (int, error)

		// Get returns Delivery details by id from data store.
		Get(id string) (*Delivery, error)

		// GetByMarketID returns Delivery details by market id from data store.
		GetByMarketID(marketID string) (*Delivery, error)

		// Create persists a new Delivery to data store.
		Create(*Delivery) error

		// Update save changes of Delivery to data store.
		Update(*Delivery) error

		// ToVerify returns a list of deliveries to process from data store.
		ToVerify(opts FindOpts) ([]Delivery, error)
	}
)

// String returns text value of a delivery status.
func (s DeliveryStatus) String() string {
	t, ok := deliveryStatusTexts[s]
	if !ok {
		return strconv.Itoa(int(s))
	}

	return t
}

// CheckCreate validates field on creating new delivery.
func (d Delivery) CheckCreate() error {
	// Check required fields.
	if err := validator.Struct(d); err != nil {
		return err
	}

	return nil
}

func (d Delivery) IsGiftOpened() *Delivery {
	if d.Assets == nil {
		return &d
	}

	opened := true
	for _, aa := range d.Assets {
		if aa.StillWrapped() {
			opened = false
			break
		}
	}

	d.GiftOpened = &opened
	return &d
}

// AddAssets handles addition of assets and remove duplicates.
func (d Delivery) AddAssets(sa []SteamAsset) *Delivery {
	d.Assets = append(d.Assets, sa...)

	keys := make(map[string]struct{})
	var unique []SteamAsset
	for _, aa := range d.Assets {
		if _, ok := keys[aa.AssetID]; ok {
			continue
		}

		keys[aa.AssetID] = struct{}{}
		unique = append(unique, aa)
	}

	d.Assets = unique
	return &d
}

// RetriesExceeded when it reached DeliveryRetryLimit reties.
func (d Delivery) RetriesExceeded() bool {
	return d.Retries > DeliveryRetryLimit
}
