package dotagiftx

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/kudarap/dotagiftx/logging"
)

const defaultCurrency = "USD"

// Market error types.
const (
	MarketErrNotFound Errors = iota + marketErrorIndex
	MarketErrRequiredID
	MarketErrRequiredFields
	MarketErrInvalidStatus
	MarketErrNotesLimit
	MarketErrInvalidPrice
	MarketErrQtyLimitPerUser
	MarketErrRequiredPartnerURL
	MarketErrInvalidBidPrice
	MarketErrInvalidAskPrice
)

// sets error text definition.
func init() {
	appErrorText[MarketErrNotFound] = "market not found"
	appErrorText[MarketErrRequiredID] = "market id is required"
	appErrorText[MarketErrRequiredFields] = "market fields are required"
	appErrorText[MarketErrInvalidStatus] = "market status not allowed"
	appErrorText[MarketErrNotesLimit] = "market notes text limit reached"
	appErrorText[MarketErrInvalidPrice] = "market price is invalid"
	appErrorText[MarketErrQtyLimitPerUser] = "market quantity limit(5) per item reached"
	appErrorText[MarketErrRequiredPartnerURL] = "market partner steam url is required"
	appErrorText[MarketErrInvalidBidPrice] = "market bid should be lower than lowest ask price"
	appErrorText[MarketErrInvalidAskPrice] = "market ask should be higher than highest bid price"
}

const (
	maxMarketNotesLen               = 200
	MaxMarketQtyLimitPerFreeUser    = 1
	MaxMarketQtyLimitPerPremiumUser = 5

	MarketAskExpirationDays = 30
	MarketBidExpirationDays = 7

	MarketSweepExpiredDays = 30
	MarketSweepRemovedDays = 60
)

// Market types.
const (
	MarketTypeAsk MarketType = 10 // default
	MarketTypeBid MarketType = 20
)

// Market statuses.
const (
	MarketStatusPending      MarketStatus = 100
	MarketStatusLive         MarketStatus = 200
	MarketStatusReserved     MarketStatus = 300
	MarketStatusSold         MarketStatus = 400
	MarketStatusBidCompleted MarketStatus = 410
	MarketStatusRemoved      MarketStatus = 500
	MarketStatusCancelled    MarketStatus = 600
	MarketStatusExpired      MarketStatus = 700
)

// Market trending score rates.
const (
	TrendScoreRateView        = 0.05
	TrendScoreRateMarketEntry = 0.01
	TrendScoreRateReserved    = 4
	TrendScoreRateSold        = 4
	TrendScoreRateBid         = 2
)

type (
	// MarketType represents market type.
	MarketType uint

	// MarketStatus represents market status.
	MarketStatus uint

	// Market represents market information.
	Market struct {
		ID             string       `json:"id"               db:"id,omitempty"`
		UserID         string       `json:"user_id"          db:"user_id,omitempty,indexed"   valid:"required"`
		ItemID         string       `json:"item_id"          db:"item_id,omitempty,indexed"   valid:"required"`
		Type           MarketType   `json:"type"             db:"type,omitempty,indexed"      valid:"required"`
		Status         MarketStatus `json:"status"           db:"status,omitempty,indexed"    valid:"required"`
		Price          float64      `json:"price"            db:"price,omitempty,indexed"     valid:"required"`
		Currency       string       `json:"currency"         db:"currency,omitempty"`
		PartnerSteamID string       `json:"partner_steam_id" db:"partner_steam_id,indexed,omitempty"`
		Notes          string       `json:"notes"            db:"notes,omitempty"`
		CreatedAt      *time.Time   `json:"created_at"       db:"created_at,omitempty,indexed"`
		UpdatedAt      *time.Time   `json:"updated_at"       db:"updated_at,omitempty,indexed"`

		InventoryStatus InventoryStatus `json:"inventory_status" db:"inventory_status,omitempty,indexed"`
		DeliveryStatus  DeliveryStatus  `json:"delivery_status"  db:"delivery_status,omitempty,indexed"`

		// Include related fields.
		User      *User      `json:"user,omitempty"      db:"user,omitempty"`
		Item      *Item      `json:"item,omitempty"      db:"item,omitempty"`
		Delivery  *Delivery  `json:"delivery,omitempty"  db:"delivery,omitempty"`
		Inventory *Inventory `json:"inventory,omitempty" db:"inventory,omitempty"`

		// reselling details.
		Resell        *bool  `json:"resell"          db:"resell,omitempty"`
		SellerSteamID string `json:"seller_steam_id" db:"seller_steam_id,omitempty"`

		// Search Indexing.
		SearchText    string `json:"-"               db:"search_text,omitempty,indexed"`
		UserRankScore int    `json:"user_rank_score" db:"user_rank_score,omitempty,indexed"`
	}

	// MarketService provides access to market service.
	MarketService interface {
		// Markets returns a list of markets.
		Markets(ctx context.Context, opts FindOpts) ([]Market, *FindMetadata, error)

		// Market returns market details by id.
		Market(ctx context.Context, id string) (*Market, error)

		// Create saves new market details.
		Create(context.Context, *Market) error

		// Update saves market details changes.
		Update(context.Context, *Market) error

		// UpdateUserRankScore sets new user ranking score on all live markets by user id.
		UpdateUserRankScore(userID string) error

		// Index composes market data for faster search and retrieval.
		// Index(ctx context.Context, id string) (*Market, error)

		// AutoCompleteBid detects if there's a matching reservation on buy order and automatically
		// resolve it by setting complete-bid status.
		AutoCompleteBid(ctx context.Context, ask Market, partnerSteamID string) error

		// Catalog returns a list of catalogs.
		Catalog(opts FindOpts) ([]Catalog, *FindMetadata, error)

		// CatalogDetails returns catalog details by item id.
		CatalogDetails(id string, opts FindOpts) (*Catalog, error)

		// TrendingCatalog returns a top 10 trending catalogs.
		TrendingCatalog(opts FindOpts) ([]Catalog, *FindMetadata, error)
	}

	// MarketStorage defines operation for market records.
	MarketStorage interface {
		// Find returns a list of markets from data store.
		Find(opts FindOpts) ([]Market, error)

		// Count returns number of market from data store.
		Count(FindOpts) (int, error)

		// Get returns a market details by id from data store.
		Get(id string) (*Market, error)

		// Create persists a new market to data store.
		Create(*Market) error

		// Update persists market changes to data store.
		Update(*Market) error

		// BaseUpdate persists market changes to data store and
		// will not update updated_at field.
		BaseUpdate(*Market) error

		// PendingInventoryStatus returns market entries that is pending for checking
		// inventory status or needs re-processing of re-process error status.
		PendingInventoryStatus(o FindOpts) ([]Market, error)

		// PendingDeliveryStatus returns market entries that is pending for checking
		// delivery status or needs re-processing of re-process error status.
		PendingDeliveryStatus(o FindOpts) ([]Market, error)

		RevalidateDeliveryStatus(o FindOpts) ([]Market, error)

		// Index composes market data for faster search and retrieval.
		Index(id string) (*Market, error)

		// UpdateUserScore sets new rank score value of all live markets by user ID.
		UpdateUserScore(userID string, rankScore int) error

		// UpdateExpiring sets live items to expired status by expiration time.
		UpdateExpiring(t MarketType, b UserBoon, expiration time.Time) (itemIDs []string, err error)

		BulkDeleteByStatus(ms MarketStatus, cutOff time.Time, limit int) error

		UpdateExpiringResell(b UserBoon) (itemIDs []string, err error)
	}
)

var MarketStatusTexts = map[MarketStatus]string{
	MarketStatusPending:      "pending",
	MarketStatusLive:         "live",
	MarketStatusReserved:     "reserved",
	MarketStatusSold:         "sold",
	MarketStatusBidCompleted: "completed",
	MarketStatusRemoved:      "removed",
	MarketStatusCancelled:    "cancelled",
	MarketStatusExpired:      "expired",
}

// CheckCreate validates field on creating new market.
func (m Market) CheckCreate() error {
	// Check required fields.
	if err := validator.Struct(m); err != nil {
		return err
	}

	// Check valid market price.
	if m.Price <= 0 {
		return MarketErrInvalidPrice
	}

	// Check market notes length.
	if len(m.Notes) > maxMarketNotesLen {
		return MarketErrNotesLimit
	}

	return nil
}

// CheckUpdate validates field on updating market.
func (m Market) CheckUpdate() error {
	if m.Notes != "" && len(m.Notes) > maxMarketNotesLen {
		return MarketErrNotesLimit
	}

	_, ok := MarketStatusTexts[m.Status]
	if m.Status != 0 && !ok {
		return MarketErrInvalidStatus
	}

	if m.Status == MarketStatusReserved && m.PartnerSteamID == "" {
		return MarketErrRequiredPartnerURL
	}

	return nil
}

// SetDefaults sets default values for a new market.
func (m Market) SetDefaults() *Market {
	m.Status = MarketStatusLive
	m.Currency = defaultCurrency
	m.Price = priceToTenths(m.Price)
	if m.Type == 0 {
		m.Type = MarketTypeAsk
	}
	return &m
}

// IsResell check if the market is a re-sell item.
func (m Market) IsResell() bool {
	return m.Type == MarketTypeAsk && m.Resell != nil && *m.Resell
}

// String returns text value of a market status.
func (s MarketStatus) String() string {
	t, ok := MarketStatusTexts[s]
	if !ok {
		return strconv.Itoa(int(s))
	}

	return t
}

// NewMarketService returns new Market service.
func NewMarketService(
	ss MarketStorage,
	us UserStorage,
	is ItemStorage,
	ts TrackStorage,
	cs CatalogStorage,
	st StatsStorage,
	vd DeliveryService,
	vi InventoryService,
	sc SteamClient,
	tp taskProcessor,
	lg logging.Logger,
) MarketService {
	return &marketService{
		ss, us,
		is,
		ts,
		cs,
		st,
		vd,
		vi,
		sc,
		tp,
		lg,
	}
}

type marketService struct {
	marketStg    MarketStorage
	userStg      UserStorage
	itemStg      ItemStorage
	trackStg     TrackStorage
	catalogStg   CatalogStorage
	statsStg     StatsStorage
	deliverySvc  DeliveryService
	inventorySvc InventoryService
	steam        SteamClient
	taskProc     taskProcessor
	logger       logging.Logger
}

func (s *marketService) Markets(ctx context.Context, opts FindOpts) ([]Market, *FindMetadata, error) {
	// Set market owner result.
	if au := AuthFromContext(ctx); au != nil {
		opts.UserID = au.UserID
	}

	res, err := s.marketStg.Find(opts)
	if err != nil {
		return nil, nil, err
	}

	if !opts.WithMeta {
		return res, nil, err
	}

	// Get total count for metadata.
	tc, err := s.marketStg.Count(opts)
	if err != nil {
		return nil, nil, err
	}

	return res, &FindMetadata{
		ResultCount: len(res),
		TotalCount:  tc,
	}, nil
}

func (s *marketService) Market(ctx context.Context, id string) (*Market, error) {
	mkt, err := s.marketStg.Get(id)
	if err != nil {
		return nil, err
	}

	// Check market ownership.
	if au := AuthFromContext(ctx); au != nil && mkt.UserID != au.UserID {
		return nil, MarketErrNotFound
	}

	return mkt, nil
}

func (s *marketService) Create(ctx context.Context, market *Market) error {
	// Set market ownership.
	au := AuthFromContext(ctx)
	if au == nil {
		return AuthErrNoAccess
	}
	market.UserID = au.UserID
	// Prevents access to create a new market when an account is flagged.
	if err := s.checkFlaggedUser(au.UserID); err != nil {
		return err
	}

	*market = *market.SetDefaults()
	if err := market.CheckCreate(); err != nil {
		return err
	}

	// Check Item existence.
	item, _ := s.itemStg.Get(market.ItemID)
	if item == nil || !item.IsActive() {
		return ItemErrNotFound
	}
	market.ItemID = item.ID

	// Check market details by type.
	switch market.Type {
	case MarketTypeAsk:
		if err := s.checkAskType(market); err != nil {
			return err
		}

		m, err := s.processShopkeepersContract(market)
		if err != nil {
			return err
		}
		market = m
	case MarketTypeBid:
		if err := s.checkBidType(market); err != nil {
			return fmt.Errorf("could not check bid type: %s", err)
		}
	}

	if err := s.marketStg.Create(market); err != nil {
		return err
	}

	bench(s.logger, "market create :: UpdateUserRankScore", func() {
		if err := s.UpdateUserRankScore(market.UserID); err != nil {
			s.logger.Errorf("could not update user rank %s: %s", market.UserID, err)
		}
	})
	bench(s.logger, "market create :: marketStg.Index", func() {
		if _, err := s.marketStg.Index(market.ID); err != nil {
			s.logger.Errorf("could not index market %s: %s", market.ItemID, err)
		}
	})
	bench(s.logger, "market create :: catalogStg.Index", func() {
		if _, err := s.catalogStg.Index(market.ItemID); err != nil {
			s.logger.Errorf("could not index item %s: %s", market.ItemID, err)
		}
	})

	// Queueing tasks for verifying post to prepare task payload.
	if market.Type == MarketTypeAsk {
		user, err := s.userStg.Get(market.UserID)
		if err != nil {
			return err
		}

		market.User = user
		market.Item = item

		// Resells should not verify items.
		if !market.IsResell() {
			if _, err = s.taskProc.Queue(ctx, user.TaskPriorityQueue(), TaskTypeVerifyInventory, market); err != nil {
				s.logger.Errorf("could not queue task: market id %s: %s", market.ID, err)
			}
		}
	}

	return nil
}

func (s *marketService) Update(ctx context.Context, market *Market) error {
	cur, err := s.checkOwnership(ctx, market.ID)
	if err != nil {
		return err
	}
	// Prevents access to update existing market when account is flagged.
	if err = s.checkFlaggedUser(cur.UserID); err != nil {
		return err
	}

	if err = market.CheckUpdate(); err != nil {
		return err
	}

	// Resolves steam profile URL input as partner steam id.
	if strings.TrimSpace(market.PartnerSteamID) != "" {
		market.PartnerSteamID, err = s.steam.ResolveVanityURL(market.PartnerSteamID)
		if err != nil {
			return err
		}

		u, err := s.userStg.Get(cur.UserID)
		if err != nil {
			return err
		}
		if u.SteamID == market.PartnerSteamID {
			return fmt.Errorf("delivering items to own account not allowed")
		}
	}
	// Try to find a matching bid and set its status to complete.
	if cur.Type == MarketTypeAsk && market.Status == MarketStatusReserved {
		if err = s.AutoCompleteBid(ctx, *cur, market.PartnerSteamID); err != nil {
			return err
		}
	}

	// Append note to existing notes.
	market.Notes = strings.TrimSpace(fmt.Sprintf("%s\n%s", cur.Notes, market.Notes))

	// Do not allow update on these fields.
	market.UserID = ""
	market.ItemID = ""
	market.Price = 0
	market.Currency = ""
	if err = s.marketStg.Update(market); err != nil {
		return err
	}

	// Queueing tasks for verifications on inventory and delivery to prepare task payload.
	if market.Type == MarketTypeAsk {
		user, err := s.userStg.Get(market.UserID)
		if err != nil {
			return err
		}
		item, err := s.itemStg.Get(market.ItemID)
		if err != nil {
			return err
		}

		market.User = user
		market.Item = item
		priority := user.TaskPriorityQueue()
		switch market.Status {
		case MarketStatusReserved:
			// Resells should not verify items.
			if !market.IsResell() {
				if _, err = s.taskProc.Queue(ctx, priority, TaskTypeVerifyInventory, market); err != nil {
					s.logger.Errorf("could not queue task: market id %s: %s", market.ID, err)
				}
			}
		case MarketStatusSold:
			if _, err = s.taskProc.Queue(ctx, priority, TaskTypeVerifyDelivery, market); err != nil {
				s.logger.Errorf("could not queue task: market id %s: %s", market.ID, err)
			}
		}
	}

	bench(s.logger, "market update :: UpdateUserRankScore", func() {
		if err = s.UpdateUserRankScore(market.UserID); err != nil {
			s.logger.Errorf("could not update user rank %s: %s", market.UserID, err)
		}
	})
	bench(s.logger, "market update :: marketStg.Index", func() {
		if _, err = s.marketStg.Index(market.ID); err != nil {
			s.logger.Errorf("could not index market %s: %s", market.ItemID, err)
		}
	})
	bench(s.logger, "market update :: catalogStg.Index", func() {
		if _, err = s.catalogStg.Index(market.ItemID); err != nil {
			s.logger.Errorf("could not index item %s: %s", market.ItemID, err)
		}
	})

	return nil
}

func (s *marketService) UpdateUserRankScore(userID string) error {
	stats, err := s.statsStg.CountUserMarketStatus(userID)
	if err != nil {
		return fmt.Errorf("error getting user market stats: %s", err)
	}

	benchS := time.Now()
	u := &User{ID: userID, MarketStats: *stats}
	u = u.CalcRankScore(*stats)
	if err = s.marketStg.UpdateUserScore(u.ID, u.RankScore); err != nil {
		return err
	}
	s.logger.Println("service/market UpdateUserScore", time.Since(benchS))
	return s.userStg.BaseUpdate(u)
}

// AutoCompleteBid detects if there's a matching reservation on buy order and automatically resolve it by setting
// complete-bid status.
func (s *marketService) AutoCompleteBid(_ context.Context, ask Market, partnerSteamID string) error {
	if ask.ItemID == "" || ask.UserID == "" || partnerSteamID == "" {
		return fmt.Errorf("ask market item id, user id, and partner steam id are required")
	}

	// Use buyer ID to get the matching market.
	buyer, err := s.userStg.Get(partnerSteamID)
	if err != nil {
		return nil
	}

	// Find matching bid market to update status.
	fo := FindOpts{
		IndexKey: "user_id",
		Filter: Market{
			Type:   MarketTypeBid,
			Status: MarketStatusLive,
			ItemID: ask.ItemID,
			UserID: buyer.ID,
		},
	}
	bids, _ := s.marketStg.Find(fo)
	if len(bids) == 0 {
		return nil
	}

	// Set complete status and seller steam id on matching bid.
	seller, err := s.userStg.Get(ask.UserID)
	if err != nil {
		return err
	}
	b := bids[0]
	b.Status = MarketStatusBidCompleted
	b.PartnerSteamID = seller.SteamID
	return s.marketStg.Update(&b)
}

func (s *marketService) checkFlaggedUser(userID string) error {
	u, err := s.userStg.Get(userID)
	if err != nil {
		return err
	}
	if err = u.CheckStatus(); err != nil {
		return err
	}

	return nil
}

func (s *marketService) processShopkeepersContract(m *Market) (*Market, error) {
	user, err := s.userStg.Get(m.UserID)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(m.SellerSteamID) == "" {
		return m, nil
	}
	if !user.HasBoon(BoonShopKeepersContract) {
		return nil, fmt.Errorf("could not find BoonShopKeepersContract")
	}

	ssid, err := s.steam.ResolveVanityURL(m.SellerSteamID)
	if err != nil {
		return nil, err
	}
	truePtr := true
	m.Resell = &truePtr
	m.SellerSteamID = ssid
	m.InventoryStatus = InventoryStatusVerified // Override verification by reseller
	return m, nil
}

func (s *marketService) checkAskType(ask *Market) error {
	user, err := s.userStg.Get(ask.UserID)
	if err != nil {
		return err
	}
	qtyLimit := MaxMarketQtyLimitPerFreeUser
	if user.HasBoon(BoonRefresherOrb) {
		qtyLimit = MaxMarketQtyLimitPerPremiumUser
	}

	if ask.PartnerSteamID != "" && user.SteamID == ask.PartnerSteamID {
		return fmt.Errorf("pairing to own account not allowed")
	}

	// Check Item max offer limit.
	qty, err := s.marketStg.Count(FindOpts{
		IndexKey: "user_id",
		Filter: Market{
			UserID: ask.UserID,
			ItemID: ask.ItemID,
			Type:   MarketTypeAsk,
			Status: MarketStatusLive,
		},
	})
	if err != nil {
		return err
	}
	if qty >= qtyLimit {
		return fmt.Errorf("market quantity limit(%d) per item reached", qtyLimit)
	}

	return nil
}

func (s *marketService) checkBidType(bid *Market) error {
	// Remove existing buy order if exists.
	res, err := s.marketStg.Find(FindOpts{
		IndexKey: "user_id",
		Filter: Market{
			UserID: bid.UserID,
			ItemID: bid.ItemID,
			Type:   MarketTypeBid,
			Status: MarketStatusLive,
		},
	})
	if err != nil {
		return err
	}
	for _, m := range res {
		m.Status = MarketStatusRemoved
		if err = s.marketStg.Update(&m); err != nil {
			return err
		}
	}

	return nil
}

func (s *marketService) checkOwnership(ctx context.Context, id string) (*Market, error) {
	au := AuthFromContext(ctx)
	if au == nil {
		return nil, AuthErrNoAccess
	}

	mkt, err := s.userMarket(au.UserID, id)
	if err != nil {
		return nil, NewXError(AuthErrForbidden, err)
	}

	if mkt == nil {
		return nil, NewXError(AuthErrForbidden, MarketErrNotFound)
	}

	return mkt, nil
}

func (s *marketService) userMarket(userID, id string) (*Market, error) {
	cur, err := s.marketStg.Get(id)
	if err != nil {
		return nil, err
	}
	if cur.UserID != userID {
		return nil, MarketErrNotFound
	}

	return cur, nil
}

func (s *marketService) Catalog(opts FindOpts) ([]Catalog, *FindMetadata, error) {
	opts.Keyword = strings.ReplaceAll(opts.Keyword, `\`, "")

	res, err := s.catalogStg.Find(opts)
	if err != nil {
		return nil, nil, err
	}

	if !opts.WithMeta {
		return res, nil, err
	}

	// Get a result and total count for metadata.
	tc, err := s.catalogStg.Count(opts)
	if err != nil {
		return nil, nil, err
	}

	return res, &FindMetadata{
		ResultCount: len(res),
		TotalCount:  tc,
	}, nil
}

func (s *marketService) TrendingCatalog(opts FindOpts) ([]Catalog, *FindMetadata, error) {
	res, err := s.catalogStg.Trending()
	if err != nil {
		return nil, nil, err
	}

	if !opts.WithMeta {
		return res, nil, err
	}

	return res, &FindMetadata{
		ResultCount: len(res),
		TotalCount:  10, // Fixed value of top 10
	}, nil
}

func (s *marketService) CatalogDetails(slug string, opts FindOpts) (*Catalog, error) {
	if slug == "" {
		return nil, CatalogErrNotFound
	}

	catalog, err := s.catalogStg.Get(slug)
	if errors.Is(err, CatalogErrNotFound) {
		i, err := s.itemStg.GetBySlug(slug)
		if err != nil {
			return nil, err
		}

		c := i.ToCatalog()
		return &c, nil

	} else if err != nil {
		return nil, err
	}

	// Override filter to specific item id.
	filter := opts.Filter.(*Market)
	filter.ItemID = catalog.ID
	opts.Filter = filter
	res, meta, err := s.Markets(context.Background(), opts)
	if err != nil {
		return nil, err
	}
	catalog.Asks = res
	catalog.Quantity = meta.TotalCount

	return catalog, err
}

type taskProcessor interface {
	Queue(ctx context.Context, p TaskPriority, t TaskType, payload interface{}) (id string, err error)
}

func bench(l logging.Logger, name string, fn func()) {
	s := time.Now()
	fn()
	l.Println("BENCH service/market", name, time.Since(s))
}

func priceToTenths(n float64) float64 {
	const dec = 100
	return math.Round(n*dec) / dec
}
