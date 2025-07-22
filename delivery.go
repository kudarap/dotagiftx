package dotagiftx

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"
)

// sets error text definition.
func init() {
	appErrorText[DeliveryErrNotFound] = "report not found"
	appErrorText[DeliveryErrRequiredID] = "report id is required"
	appErrorText[DeliveryErrRequiredFields] = "report fields are required"

	appErrorText[InventoryErrNotFound] = "report not found"
	appErrorText[InventoryErrRequiredID] = "report id is required"
	appErrorText[InventoryErrRequiredFields] = "report fields are required"
}

// Delivery error types.
const (
	DeliveryErrNotFound Errors = iota + deliveryErrorIndex
	DeliveryErrRequiredID
	DeliveryErrRequiredFields
)

// Inventory error types.
const (
	InventoryErrNotFound Errors = iota + inventoryErrorIndex
	InventoryErrRequiredID
	InventoryErrRequiredFields
)

// DeliveryRetryLimit max retry to process verification.
const DeliveryRetryLimit = 30

// Delivery statuses.
const (
	// DeliveryStatusNoHit buyer's inventory successfully parsed,
	// but the item did not find any in match.
	DeliveryStatusNoHit DeliveryStatus = 100

	// DeliveryStatusNameVerified item exists on buyer's inventory base on the item name challenge.
	//
	// No-gift info might mean:
	// 1. Buyer cleared the gift information
	// 2. Buyer is the original owner of item
	// 3. Item might come from another source
	DeliveryStatusNameVerified DeliveryStatus = 200

	// DeliveryStatusSenderVerified both item existence and gift information matched the seller's avatar name. We could
	// also use the date received to check against delivery data to strengthen its validity.
	DeliveryStatusSenderVerified DeliveryStatus = 300

	// DeliveryStatusPrivate buyer's inventory is not visible to the public, and we can do nothing about it.
	DeliveryStatusPrivate DeliveryStatus = 400

	// DeliveryStatusError error occurred during API request or parsing inventory error.
	DeliveryStatusError DeliveryStatus = 500
)

// Inventory statuses.
const (
	// InventoryStatusNoHit buyer's inventory successfully parsed, but the item did not find any in match.
	InventoryStatusNoHit InventoryStatus = 100

	// InventoryStatusVerified item exists on inventory base on the item name challenge.
	InventoryStatusVerified InventoryStatus = 200

	// InventoryStatusPrivate buyer's inventory is not visible to the public, and we can do nothing about it.
	InventoryStatusPrivate InventoryStatus = 400

	// InventoryStatusError error occurred during API request or parsing inventory error.
	InventoryStatusError InventoryStatus = 500
)

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
		VerifiedBy       string         `json:"verified_by"        db:"verified_by,omitempty,indexed"`
		ElapsedMs        int64          `json:"elapsed_ms"         db:"elapsed_ms,omitempty,indexed"`
		CreatedAt        *time.Time     `json:"created_at"         db:"created_at,omitempty,indexed,omitempty"`
		UpdatedAt        *time.Time     `json:"updated_at"         db:"updated_at,omitempty,indexed,omitempty"`
	}

	// DeliveryService provides access to Delivery service.
	DeliveryService interface {
		// Deliveries return a list of deliveries.
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

	// InventoryStatus represents inventory status.
	InventoryStatus uint

	// Inventory represents steam inventory.
	Inventory struct {
		ID          string          `json:"id"           db:"id,omitempty,omitempty"`
		MarketID    string          `json:"market_id"    db:"market_id,omitempty,indexed" valid:"required"`
		Status      InventoryStatus `json:"status"       db:"status,omitempty,indexed"    valid:"required"`
		Assets      []SteamAsset    `json:"steam_assets" db:"steam_assets,omitempty"`
		Retries     int             `json:"retries"      db:"retries,omitempty"`
		BundleCount int             `json:"bundle_count" db:"bundle_count,omitempty"`
		VerifiedBy  string          `json:"verified_by"  db:"verified_by,omitempty,indexed"`
		ElapsedMs   int64           `json:"elapsed_ms"   db:"elapsed_ms,omitempty,indexed"`
		CreatedAt   *time.Time      `json:"created_at"   db:"created_at,omitempty,indexed,omitempty"`
		UpdatedAt   *time.Time      `json:"updated_at"   db:"updated_at,omitempty,indexed,omitempty"`
	}

	// InventoryService provides access to Inventory service.
	InventoryService interface {
		// Inventories returns a list of deliveries.
		Inventories(opts FindOpts) ([]Inventory, *FindMetadata, error)

		// Inventory returns Inventory details by id.
		Inventory(id string) (*Inventory, error)

		// Set saves new Inventory details.
		Set(context.Context, *Inventory) error
	}

	// InventoryStorage defines operation for Inventory records.
	InventoryStorage interface {
		// Find returns a list of inventories from data store.
		Find(opts FindOpts) ([]Inventory, error)

		// Count returns number of inventories from data store.
		Count(FindOpts) (int, error)

		// Get returns an Inventory details by id from data store.
		Get(id string) (*Inventory, error)

		// GetByMarketID returns Inventory details by market id from data store.
		GetByMarketID(marketID string) (*Inventory, error)

		// Create persists a new Inventory to data store.
		Create(*Inventory) error

		// Update save changes of Inventory to data store.
		Update(*Inventory) error
	}
)

// CheckCreate validates field on creating new delivery.
func (d Delivery) CheckCreate() error {
	// Check the required fields.
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

// CheckCreate validates field on creating new inventory.
func (i Inventory) CheckCreate() error {
	// Check the required fields.
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

// RetriesExceeded when it reached 5 retries.
func (i Inventory) RetriesExceeded() bool {
	return i.Retries > 3
}

var deliveryStatusTexts = map[DeliveryStatus]string{
	DeliveryStatusNoHit:          "no hit",
	DeliveryStatusNameVerified:   "name verified",
	DeliveryStatusSenderVerified: "sender verified",
	DeliveryStatusPrivate:        "private",
	DeliveryStatusError:          "error",
}

// String returns text value of a delivery status.
func (s DeliveryStatus) String() string {
	t, ok := deliveryStatusTexts[s]
	if !ok {
		return strconv.Itoa(int(s))
	}

	return t
}

var inventoryStatusTexts = map[InventoryStatus]string{
	InventoryStatusNoHit:    "no hit",
	InventoryStatusVerified: "verified",
	InventoryStatusPrivate:  "private",
	InventoryStatusError:    "error",
}

// String returns text value of an inventory status.
func (s InventoryStatus) String() string {
	t, ok := inventoryStatusTexts[s]
	if !ok {
		return strconv.Itoa(int(s))
	}

	return t
}

// NewDeliveryService returns a new delivery service.
func NewDeliveryService(rs DeliveryStorage, ms MarketStorage) DeliveryService {
	return &deliveryService{rs, ms}
}

type deliveryService struct {
	deliveryStg DeliveryStorage
	marketStg   MarketStorage
}

func (s *deliveryService) Deliveries(opts FindOpts) ([]Delivery, *FindMetadata, error) {
	res, err := s.deliveryStg.Find(opts)
	if err != nil {
		return nil, nil, err
	}

	if !opts.WithMeta {
		return res, nil, err
	}

	// Get a result and total count for metadata.
	tc, err := s.deliveryStg.Count(opts)
	if err != nil {
		return nil, nil, err
	}

	return res, &FindMetadata{
		ResultCount: len(res),
		TotalCount:  tc,
	}, nil
}

func (s *deliveryService) Delivery(id string) (*Delivery, error) {
	inv, err := s.deliveryStg.Get(id)
	if err != nil && !errors.Is(err, DeliveryErrNotFound) {
		return nil, err
	}
	if inv != nil {
		return inv, nil
	}

	// If we can't find using id, let's try market id
	return s.deliveryStg.GetByMarketID(id)
}

func (s *deliveryService) DeliveryByMarketID(marketID string) (*Delivery, error) {
	return s.deliveryStg.GetByMarketID(marketID)
}

func (s *deliveryService) Set(_ context.Context, del *Delivery) error {
	if err := del.CheckCreate(); err != nil {
		return NewXError(DeliveryErrRequiredFields, err)
	}

	defer func() {
		if _, err := s.marketStg.Index(del.MarketID); err != nil {
			log.Printf("could not index market %s: %s", del.MarketID, err)
		}
	}()

	// Detect if there is still unopened gift.
	del = del.IsGiftOpened()

	// Update market delivery status.
	if err := s.marketStg.BaseUpdate(&Market{
		ID:             del.MarketID,
		DeliveryStatus: del.Status,
	}); err != nil {
		return err
	}

	cur, _ := s.DeliveryByMarketID(del.MarketID)
	if cur != nil {
		del.ID = cur.ID
		del.Retries = cur.Retries + 1
		del = del.AddAssets(cur.Assets)
		return s.deliveryStg.Update(del)
	}

	return s.deliveryStg.Create(del)
}

// NewInventoryService returns new inventory service.
func NewInventoryService(rs InventoryStorage, ms MarketStorage, cs CatalogStorage) InventoryService {
	return &inventoryService{rs, ms, cs}
}

type inventoryService struct {
	inventoryStg InventoryStorage
	marketStg    MarketStorage
	catalogStg   CatalogStorage
}

func (s *inventoryService) Inventories(opts FindOpts) ([]Inventory, *FindMetadata, error) {
	res, err := s.inventoryStg.Find(opts)
	if err != nil {
		return nil, nil, err
	}

	if !opts.WithMeta {
		return res, nil, err
	}

	// Get a result and total count for metadata.
	tc, err := s.inventoryStg.Count(opts)
	if err != nil {
		return nil, nil, err
	}

	return res, &FindMetadata{
		ResultCount: len(res),
		TotalCount:  tc,
	}, nil
}

func (s *inventoryService) Inventory(id string) (*Inventory, error) {
	inv, err := s.inventoryStg.Get(id)
	if err != nil && !errors.Is(err, InventoryErrNotFound) {
		return nil, err
	}
	if inv != nil {
		return inv, nil
	}

	// If we can't find using id, let's try market id
	return s.inventoryStg.GetByMarketID(id)
}

func (s *inventoryService) InventoryByMarketID(marketID string) (*Inventory, error) {
	return s.inventoryStg.GetByMarketID(marketID)
}

func (s *inventoryService) Set(_ context.Context, inv *Inventory) error {
	if err := inv.CheckCreate(); err != nil {
		return NewXError(InventoryErrRequiredFields, err)
	}

	defer func() {
		mkt, err := s.marketStg.Index(inv.MarketID)
		if err != nil {
			log.Printf("could not index market %s: %s", inv.MarketID, err)
		}
		if _, err = s.catalogStg.Index(mkt.ItemID); err != nil {
			log.Printf("could not index catalog %s: %s", inv.MarketID, err)
		}
	}()

	// Update market Inventory status.
	if err := s.marketStg.BaseUpdate(&Market{
		ID:              inv.MarketID,
		InventoryStatus: inv.Status,
	}); err != nil {
		return err
	}

	// Process bundle count.
	inv.BundleCount = inv.CountBundles()
	cur, _ := s.InventoryByMarketID(inv.MarketID)
	if cur != nil {
		inv.ID = cur.ID
		inv.Retries = cur.Retries + 1
		return s.inventoryStg.Update(inv)
	}

	return s.inventoryStg.Create(inv)
}
