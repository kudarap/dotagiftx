package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/errors"
	"github.com/kudarap/dotagiftx/gokit/log"
)

type TaskProcessor interface {
	Queue(ctx context.Context, p dgx.TaskPriority, t dgx.TaskType, payload interface{}) (id string, err error)
}

// NewMarket returns new Market service.
func NewMarket(
	ss dgx.MarketStorage,
	us dgx.UserStorage,
	is dgx.ItemStorage,
	ts dgx.TrackStorage,
	cs dgx.CatalogStorage,
	st dgx.StatsStorage,
	vd dgx.DeliveryService,
	vi dgx.InventoryService,
	sc dgx.SteamClient,
	tp TaskProcessor,
	lg log.Logger,
) dgx.MarketService {
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
	marketStg    dgx.MarketStorage
	userStg      dgx.UserStorage
	itemStg      dgx.ItemStorage
	trackStg     dgx.TrackStorage
	catalogStg   dgx.CatalogStorage
	statsStg     dgx.StatsStorage
	deliverySvc  dgx.DeliveryService
	inventorySvc dgx.InventoryService
	steam        dgx.SteamClient
	taskProc     TaskProcessor
	logger       log.Logger
}

func (s *marketService) Markets(ctx context.Context, opts dgx.FindOpts) ([]dgx.Market, *dgx.FindMetadata, error) {
	// Set market owner result.
	if au := dgx.AuthFromContext(ctx); au != nil {
		opts.UserID = au.UserID
	}

	res, err := s.marketStg.Find(opts)
	if err != nil {
		return nil, nil, err
	}

	// Assign inventory and delivery status.
	for i, mkt := range res {
		res[i] = mkt
	}

	if !opts.WithMeta {
		return res, nil, err
	}

	// Get total count for metadata.
	tc, err := s.marketStg.Count(opts)
	if err != nil {
		return nil, nil, err
	}

	return res, &dgx.FindMetadata{
		ResultCount: len(res),
		TotalCount:  tc,
	}, nil
}

func (s *marketService) Market(ctx context.Context, id string) (*dgx.Market, error) {
	mkt, err := s.marketStg.Get(id)
	if err != nil {
		return nil, err
	}

	// Check market ownership.
	if au := dgx.AuthFromContext(ctx); au != nil && mkt.UserID != au.UserID {
		return nil, dgx.MarketErrNotFound
	}

	return mkt, nil
}

func (s *marketService) Create(ctx context.Context, market *dgx.Market) error {
	// Set market ownership.
	au := dgx.AuthFromContext(ctx)
	if au == nil {
		return dgx.AuthErrNoAccess
	}
	market.UserID = au.UserID
	// Prevents access to create new market when account is flagged.
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
		return dgx.ItemErrNotFound
	}
	market.ItemID = item.ID

	// Check market details by type.
	switch market.Type {
	case dgx.MarketTypeAsk:
		if err := s.checkAskType(market); err != nil {
			return err
		}

		m, err := s.processShopkeepersContract(market)
		if err != nil {
			return err
		}
		market = m
	case dgx.MarketTypeBid:
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
	if market.Type == dgx.MarketTypeAsk {
		user, err := s.userStg.Get(market.UserID)
		if err != nil {
			return err
		}

		market.User = user
		market.Item = item

		// Resells should not verify items.
		if !market.IsResell() {
			if _, err = s.taskProc.Queue(ctx, user.TaskPriorityQueue(), dgx.TaskTypeVerifyInventory, market); err != nil {
				s.logger.Errorf("could not queue task: market id %s: %s", market.ID, err)
			}
		}
	}

	return nil
}

func (s *marketService) Update(ctx context.Context, market *dgx.Market) error {
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
	if cur.Type == dgx.MarketTypeAsk && market.Status == dgx.MarketStatusReserved {
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
	if market.Type == dgx.MarketTypeAsk {
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
		case dgx.MarketStatusReserved:
			// Resells should not verify items.
			if !market.IsResell() {
				if _, err = s.taskProc.Queue(ctx, priority, dgx.TaskTypeVerifyInventory, market); err != nil {
					s.logger.Errorf("could not queue task: market id %s: %s", market.ID, err)
				}
			}
		case dgx.MarketStatusSold:
			if _, err = s.taskProc.Queue(ctx, priority, dgx.TaskTypeVerifyDelivery, market); err != nil {
				s.logger.Errorf("could not queue task: market id %s: %s", market.ID, err)
			}
		}
	}

	//if err = s.UpdateUserRankScore(market.UserID); err != nil {
	//	return err
	//}
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
	u := &dgx.User{ID: userID, MarketStats: *stats}
	u = u.CalcRankScore(*stats)
	if err = s.marketStg.UpdateUserScore(u.ID, u.RankScore); err != nil {
		return err
	}
	s.logger.Println("service/market UpdateUserScore", time.Now().Sub(benchS))
	return s.userStg.BaseUpdate(u)
}

// AutoCompleteBid detects if there's matching reservation on buy order and automatically
// resolve it by setting complete-bid status.
func (s *marketService) AutoCompleteBid(_ context.Context, ask dgx.Market, partnerSteamID string) error {
	if ask.ItemID == "" || ask.UserID == "" || partnerSteamID == "" {
		return fmt.Errorf("ask market item id, user id, and partner steam id are required")
	}

	// Use buyer ID to get the matching market.
	buyer, err := s.userStg.Get(partnerSteamID)
	if err != nil {
		return nil
	}

	// Find matching bid market to update status.
	fo := dgx.FindOpts{
		Filter: dgx.Market{
			Type:   dgx.MarketTypeBid,
			Status: dgx.MarketStatusLive,
			ItemID: ask.ItemID,
			UserID: buyer.ID,
		},
	}
	bids, _ := s.marketStg.Find(fo)
	if bids == nil || len(bids) == 0 {
		return nil
	}

	// Set complete status and seller steam id on matching bid.
	seller, err := s.userStg.Get(ask.UserID)
	if err != nil {
		return err
	}
	b := bids[0]
	b.Status = dgx.MarketStatusBidCompleted
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

func (s *marketService) processShopkeepersContract(m *dgx.Market) (*dgx.Market, error) {
	user, err := s.userStg.Get(m.UserID)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(m.SellerSteamID) == "" {
		return m, nil
	}
	if !user.HasBoon(dgx.BoonShopKeepersContract) {
		return nil, fmt.Errorf("could not find BoonShopKeepersContract")
	}

	ssid, err := s.steam.ResolveVanityURL(m.SellerSteamID)
	if err != nil {
		return nil, err
	}
	truePtr := true
	m.Resell = &truePtr
	m.SellerSteamID = ssid
	m.InventoryStatus = dgx.InventoryStatusVerified // Override verification by reseller
	return m, nil
}

func (s *marketService) checkAskType(ask *dgx.Market) error {
	//if err := s.restrictMatchingPriceValue(ask); err != nil {
	//	return err
	//}

	user, err := s.userStg.Get(ask.UserID)
	if err != nil {
		return err
	}
	qtyLimit := dgx.MaxMarketQtyLimitPerFreeUser
	if user.HasBoon(dgx.BoonRefresherOrb) {
		qtyLimit = dgx.MaxMarketQtyLimitPerPremiumUser
	}

	if ask.PartnerSteamID != "" && user.SteamID == ask.PartnerSteamID {
		return fmt.Errorf("pairing to own account not allowed")
	}

	// Check Item max offer limit.
	qty, err := s.marketStg.Count(dgx.FindOpts{
		Filter: dgx.Market{
			ItemID: ask.ItemID,
			Type:   dgx.MarketTypeAsk,
			Status: dgx.MarketStatusLive,
		},
		UserID: ask.UserID,
	})
	if err != nil {
		return err
	}
	if qty >= qtyLimit {
		return fmt.Errorf("market quantity limit(%d) per item reached", qtyLimit)
	}

	return nil
}

func (s *marketService) checkBidType(bid *dgx.Market) error {
	//if err := s.restrictMatchingPriceValue(bid); err != nil {
	//	return err
	//}

	// Remove existing buy order if exists.
	res, err := s.marketStg.Find(dgx.FindOpts{
		Filter: dgx.Market{
			ItemID: bid.ItemID,
			Type:   dgx.MarketTypeBid,
			Status: dgx.MarketStatusLive,
		},
		UserID: bid.UserID,
	})
	if err != nil {
		return err
	}
	for _, m := range res {
		m.Status = dgx.MarketStatusRemoved
		if err = s.marketStg.Update(&m); err != nil {
			return err
		}
	}

	return nil
}

// restrictMatchingPriceValue restricts market price against its counter-part entry.
// 1. market bid price should lower than the lowest ask price.
// 2. market ask price should higher than the highest bid price.
// This was design to enforce the user to check available offers or orders
// with desired price value.
// Update 2021/03/08: It turns out some users are picky on which user they
// want to get the item from, which is very reasonable, and will disable this restriction for now.
func (s *marketService) restrictMatchingPriceValue(mkt *dgx.Market) error {
	switch mkt.Type {
	case dgx.MarketTypeAsk:
		bid, err := s.catalogStg.Index(mkt.ItemID)
		if err != nil {
			return err
		}
		if bid.Quantity != 0 && bid.HighestBid > mkt.Price {
			return dgx.MarketErrInvalidAskPrice
		}
	case dgx.MarketTypeBid:
		ask, err := s.catalogStg.Index(mkt.ItemID)
		if err != nil {
			return err
		}
		if ask.Quantity != 0 && ask.LowestAsk < mkt.Price {
			return dgx.MarketErrInvalidBidPrice
		}
	}

	return nil
}

func (s *marketService) checkOwnership(ctx context.Context, id string) (*dgx.Market, error) {
	au := dgx.AuthFromContext(ctx)
	if au == nil {
		return nil, dgx.AuthErrNoAccess
	}

	mkt, err := s.userMarket(au.UserID, id)
	if err != nil {
		return nil, errors.New(dgx.AuthErrForbidden, err)
	}

	if mkt == nil {
		return nil, errors.New(dgx.AuthErrForbidden, dgx.MarketErrNotFound)
	}

	return mkt, nil
}

func (s *marketService) userMarket(userID, id string) (*dgx.Market, error) {
	cur, err := s.marketStg.Get(id)
	if err != nil {
		return nil, err
	}
	if cur.UserID != userID {
		return nil, dgx.MarketErrNotFound
	}

	return cur, nil
}

func bench(l log.Logger, name string, fn func()) {
	s := time.Now()
	fn()
	l.Println("BENCH service/market", name, time.Now().Sub(s))
}
