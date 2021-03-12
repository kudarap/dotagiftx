package core

import (
	"time"
)

/// Delivery represents steam inventory delivery data.
type Delivery struct {
	ID               string         `json:"id"             db:"id,omitempty"`
	MarketID         string         `json:"market_id"`
	BuyerConfirmed   *bool          `json:"buyer_confirmed"`
	BuyerConfirmedAt *time.Time     `json:"buyer_confirmed_at"`
	Status           DeliveryStatus `json:"status"`
	Assets           []SteamAsset   `json:"steam_assets"`
	CreatedAt        *time.Time     `json:"created_at"     db:"created_at,omitempty,indexed"`
	UpdatedAt        *time.Time     `json:"updated_at"     db:"updated_at,omitempty,indexed"`
}

// DeliveryStatus represents delivery status.
type DeliveryStatus uint

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
	// public and we can do nothing about it.
	DeliveryStatusPrivate DeliveryStatus = 400

	// DeliveryStatusError error occurred during API request or
	// parsing inventory error.
	DeliveryStatusError DeliveryStatus = 500
)
