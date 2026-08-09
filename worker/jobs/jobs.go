package jobs

import (
	"context"
	"time"

	dotagiftx2 "github.com/kudarap/dotagiftx/dotagiftx"
)

type cacheRemover interface {
	BulkDel(keyPrefix string) error
}

// marketRepository provides access to market storage methods used by jobs.
type marketRepository interface {
	// Find returns a list of markets from data store.
	Find(ctx context.Context, opts dotagiftx2.FindOpts) ([]dotagiftx2.Market, error)

	// PendingInventoryStatus returns market entries that is pending for checking
	// inventory status or needs re-processing of re-process error status.
	PendingInventoryStatus(ctx context.Context, o dotagiftx2.FindOpts) ([]dotagiftx2.Market, error)

	// PendingDeliveryStatus returns market entries that is pending for checking
	// delivery status or needs re-processing of re-process error status.
	PendingDeliveryStatus(ctx context.Context, o dotagiftx2.FindOpts) ([]dotagiftx2.Market, error)

	// UpdateExpiring sets live items to expired status by expiration time.
	UpdateExpiring(ctx context.Context, t dotagiftx2.MarketType, b dotagiftx2.UserBoon, expiration time.Time) (itemIDs []string, err error)

	// BulkDeleteByStatus deletes markets by status and cutoff time.
	BulkDeleteByStatus(ctx context.Context, ms dotagiftx2.MarketStatus, cutOff time.Time, limit int) error

	// UpdateExpiringResell sets resell markets to expired status.
	UpdateExpiringResell(ctx context.Context, b dotagiftx2.UserBoon) (itemIDs []string, err error)
}

// deliveryRepository provides access to delivery storage methods used by jobs.
type deliveryRepository interface {
	// ToVerify returns a list of deliveries to process from data store.
	ToVerify(ctx context.Context, opts dotagiftx2.FindOpts) ([]dotagiftx2.Delivery, error)
}

// catalogRepository provides access to catalog storage methods used by jobs.
type catalogRepository interface {
	// Index persists a new catalog to data store.
	Index(ctx context.Context, itemID string) (*dotagiftx2.Catalog, error)
}

// userRepository provides access to user storage methods used by jobs.
type userRepository interface {
	// ExpiringSubscribers return a list of users that has expiring subscription.
	ExpiringSubscribers(ctx context.Context, now time.Time) ([]dotagiftx2.User, error)

	// PurgeSubscription removes subscription data and boons.
	PurgeSubscription(ctx context.Context, userID string) error
}

// deliveryService provides access to delivery service methods used by jobs.
type deliveryService interface {
	// Set saves new Delivery details.
	Set(context.Context, *dotagiftx2.Delivery) error
}

// inventoryService provides access to inventory service methods used by jobs.
type inventoryService interface {
	// Inventories returns a list of inventories.
	Inventories(ctx context.Context, opts dotagiftx2.FindOpts) ([]dotagiftx2.Inventory, *dotagiftx2.FindMetadata, error)

	// Set saves new Inventory details.
	Set(context.Context, *dotagiftx2.Inventory) error
}
