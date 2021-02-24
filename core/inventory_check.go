package core

import "time"

const (
	InventoryCheckTypeExists    InventoryCheckType = 10
	InventoryCheckTypeDelivered InventoryCheckType = 20

	InventoryCheckStatusPending   InventoryCheckStatus = 10
	InventoryCheckStatusError     InventoryCheckStatus = 20
	InventoryCheckStatusPrivate   InventoryCheckStatus = 30
	InventoryCheckStatusVerified  InventoryCheckStatus = 40
	InventoryCheckStatusDelivered InventoryCheckStatus = 50

	VerifiedDeliveryTypeAuto      VerifiedDeliveryType = 10 // automatically done by the server worker
	VerifiedDeliveryTypeBuyerConf VerifiedDeliveryType = 20 // manually confirmed by the buyer

	VerifiedDeliveryLevelError      VerifiedDeliveryLevel = 10 // error processing the verification
	VerifiedDeliveryLevelPrivate    VerifiedDeliveryLevel = 20 // buyer's inventory is private
	VerifiedDeliveryLevelNoHit      VerifiedDeliveryLevel = 40 // item name did not matched on buyer's inventory
	VerifiedDeliveryLevelBuyerConf  VerifiedDeliveryLevel = 50 // manually confirmed by the buyer
	VerifiedDeliveryLevelItem       VerifiedDeliveryLevel = 60 // item exist but gifter name not exist
	VerifiedDeliveryLevelItemGifted VerifiedDeliveryLevel = 70 // item exist and gifter name matched
)

type (
	InventoryCheckType   uint
	InventoryCheckStatus uint
	InventoryCheckLevel  uint

	InventoryCheck struct {
		ID        string
		MarketID  string
		Type      InventoryCheckType
		Status    InventoryCheckStatus
		AssetIDs  []InvAsset
		CreatedAt *time.Time
		UpdatedAt *time.Time
	}

	VerifiedDeliveryType  uint
	VerifiedDeliveryLevel uint
	VerifiedDelivery      struct {
		ID        string
		MarketID  string
		Level     string
		InvAssets []InvAsset
		CreatedAt *time.Time
		UpdatedAt *time.Time
	}

	InvAsset struct {
		AssetID      string   `json:"asset_id"`
		Name         string   `json:"name"`
		Type         string   `json:"type"`
		Hero         string   `json:"hero"`
		GiftFrom     string   `json:"gift_from"`
		DateReceived string   `json:"date_received"`
		Dedication   string   `json:"dedication"`
		GiftOnce     bool     `json:"gift_once"`
		NotTradable  bool     `json:"not_tradable"`
		Descriptions []string `json:"descriptions"`
	}
)
