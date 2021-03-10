package core

import (
	"context"
	"time"
)

// Verified Delivery statuses.
const (
	// VerifiedDeliveryStatusPending on-going or queued for parsing inventory.
	VerifiedDeliveryStatusPending VerifiedDeliveryStatus = 100

	// VerifiedDeliveryStatusPassed done parsing buyer inventory and pass the verification challenge.
	VerifiedDeliveryStatusPassed VerifiedDeliveryStatus = 200

	// VerifiedDeliveryStatusBuyerAck manually verified by buyer. Special used case for private
	// inventories but requires login from buyer.
	VerifiedDeliveryStatusBuyerAck VerifiedDeliveryStatus = 210

	// VerifiedDeliveryStatusFailed done parsing buyer inventory but did not pass verification challenge.
	VerifiedDeliveryStatusFailed VerifiedDeliveryStatus = 400

	// VerifiedDeliveryStatusError inventory could not be parsed or an error occurred while requesting.
	// This could mean steam api is reaching its request limit or the inventory is private.
	VerifiedDeliveryStatusError VerifiedDeliveryStatus = 500
)

type (
	VerifiedDeliveryStatus uint

	VerifiedDeliveryX struct {
		ID           string                 `json:"id"             db:"id,omitempty"`
		UserID       string                 `json:"user_id"        db:"user_id,omitempty,indexed"    valid:"required"`
		ItemID       string                 `json:"item_id"        db:"item_id,omitempty,indexed"    valid:"required"`
		BuyerSteamID string                 `json:"buyer_steam_id" db:"buyer_steam_id,omitempty"     valid:"required"`
		Status       VerifiedDeliveryStatus `json:"status"         db:"status,indexed"`
		RawData      string                 `json:"raw_data"       db:"raw_data"`
		Retries      int                    `json:"retries"        db:"retries"`
		CreatedAt    *time.Time             `json:"created_at"     db:"created_at,omitempty,indexed"`
		UpdatedAt    *time.Time             `json:"updated_at"     db:"updated_at,omitempty,indexed"`
		// Include related fields.
		User *User `json:"user,omitempty" db:"user,omitempty"`
		Item *Item `json:"item,omitempty" db:"item,omitempty"`
	}

	VerifiedService interface {
		VerifiedDeliveries(ctx context.Context, opts FindOpts) ([]VerifiedDeliveryX, *FindMetadata, error)
		VerifiedDelivery(ctx context.Context, id string) (VerifiedDeliveryX, error)
		Verify(ctx context.Context, marketID, buyerSteamID string) error
	}

	VerifiedStorage interface {
	}
)

func NewVerifiedDelivery(userID, itemID, buyerSteamID string) VerifiedDeliveryX {
	return VerifiedDeliveryX{}
}

func (d VerifiedDeliveryX) ReduceRawData() string {
	panic("implement me")
}

func (d VerifiedDeliveryX) ParseRawData() string {
	panic("implement me")
}

// Challenge validates delivery
// 1. check for gift from the seller
// 2. check item name
func (d VerifiedDeliveryX) Challenge() (passed bool) {
	panic("implement me")
}
