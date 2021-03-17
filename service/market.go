package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/errors"
	"github.com/kudarap/dotagiftx/gokit/log"
)

type Dispatcher interface {
	VerifyDelivery(marketID string)
	VerifyInventory(userID string)
}

// NewMarket returns new Market service.
func NewMarket(
	ss core.MarketStorage,
	us core.UserStorage,
	is core.ItemStorage,
	ts core.TrackStorage,
	cs core.CatalogStorage,
	vd core.DeliveryService,
	vi core.InventoryService,
	sc core.SteamClient,
	dp Dispatcher,
	lg log.Logger,
) core.MarketService {
	return &marketService{
		ss, us,
		is,
		ts,
		cs,
		vd,
		vi,
		sc,
		dp,
		lg,
	}
}

type marketService struct {
	marketStg    core.MarketStorage
	userStg      core.UserStorage
	itemStg      core.ItemStorage
	trackStg     core.TrackStorage
	catalogStg   core.CatalogStorage
	deliverySvc  core.DeliveryService
	inventorySvc core.InventoryService
	steam        core.SteamClient
	dispatch     Dispatcher
	logger       log.Logger
}

func (s *marketService) Markets(ctx context.Context, opts core.FindOpts) ([]core.Market, *core.FindMetadata, error) {
	// Set market owner result.
	if au := core.AuthFromContext(ctx); au != nil {
		opts.UserID = au.UserID
	}

	res, err := s.marketStg.Find(opts)
	if err != nil {
		return nil, nil, err
	}

	// Assign inventory and delivery status.
	for i, mkt := range res {
		s.getRelatedVerifiedStatus(&mkt)
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

	return res, &core.FindMetadata{
		ResultCount: len(res),
		TotalCount:  tc,
	}, nil
}

func (s *marketService) Market(ctx context.Context, id string) (*core.Market, error) {
	mkt, err := s.marketStg.Get(id)
	if err != nil {
		return nil, err
	}

	// Check market ownership.
	if au := core.AuthFromContext(ctx); au != nil && mkt.UserID != au.UserID {
		return nil, core.MarketErrNotFound
	}

	s.getRelatedFields(mkt)
	s.getRelatedVerifiedStatus(mkt)
	return mkt, nil
}

func (s *marketService) getRelatedFields(mkt *core.Market) {
	mkt.User, _ = s.userStg.Get(mkt.UserID)
	mkt.Item, _ = s.itemStg.Get(mkt.ItemID)
}

func (s *marketService) getRelatedVerifiedStatus(mkt *core.Market) {
	mkt.Inventory, _ = s.inventorySvc.Inventory(mkt.ID)
	mkt.Delivery, _ = s.deliverySvc.Delivery(mkt.ID)
}

func (s *marketService) Create(ctx context.Context, mkt *core.Market) error {
	// Set market ownership.
	au := core.AuthFromContext(ctx)
	if au == nil {
		return core.AuthErrNoAccess
	}
	mkt.UserID = au.UserID
	// Prevents access to create new market when account is flagged.
	if err := s.checkFlaggedUser(au.UserID); err != nil {
		return err
	}

	mkt.SetDefaults()
	if err := mkt.CheckCreate(); err != nil {
		return err
	}

	// Check Item existence.
	i, _ := s.itemStg.Get(mkt.ItemID)
	if i == nil {
		return core.ItemErrNotFound
	}
	mkt.ItemID = i.ID

	// Check market details by type.
	switch mkt.Type {
	case core.MarketTypeAsk:
		if err := s.checkAskType(mkt); err != nil {
			return err
		}
	case core.MarketTypeBid:
		if err := s.checkBidType(mkt); err != nil {
			return err
		}
	}

	if err := s.marketStg.Create(mkt); err != nil {
		return err
	}

	//go func() {
	if _, err := s.catalogStg.Index(mkt.ItemID); err != nil {
		s.logger.Errorf("could not index item %s: %s", mkt.ItemID, err)
	}
	//}()
	if mkt.Type == core.MarketTypeAsk {
		s.dispatch.VerifyInventory(mkt.UserID)
	}

	return nil
}

func (s *marketService) checkAskType(ask *core.Market) error {
	//if err := s.restrictMatchingPriceValue(ask); err != nil {
	//	return err
	//}

	// Check Item max offer limit.
	qty, err := s.marketStg.Count(core.FindOpts{
		Filter: core.Market{
			ItemID: ask.ItemID,
			Type:   core.MarketTypeAsk,
			Status: core.MarketStatusLive,
		},
		UserID: ask.UserID,
	})
	if err != nil {
		return err
	}
	if qty >= core.MaxMarketQtyLimitPerUser {
		return core.MarketErrQtyLimitPerUser
	}

	return nil
}

func (s *marketService) checkBidType(bid *core.Market) error {
	//if err := s.restrictMatchingPriceValue(bid); err != nil {
	//	return err
	//}

	// Remove existing buy order if exists.
	res, err := s.marketStg.Find(core.FindOpts{
		Filter: core.Market{
			ItemID: bid.ItemID,
			Type:   core.MarketTypeBid,
			Status: core.MarketStatusLive,
		},
		UserID: bid.UserID,
	})
	if err != nil {
		return err
	}
	for _, m := range res {
		m.Status = core.MarketStatusRemoved
		if err := s.marketStg.Update(&m); err != nil {
			return err
		}
	}

	return nil
}

// restrictMatchingPriceValue restricts market price against its counter-part entry.
// 1. market bid price should lower than lowest ask price.
// 2. market ask price should higher than highest bid price.
// This was design to enforced the user to check available offers or orders
// with desired price value.
// Update 2021/03/08: It turns out some of the user are picky on which user they
// want to get the item from, which is very reasonable, and will disable this restriction for now.
func (s *marketService) restrictMatchingPriceValue(mkt *core.Market) error {
	switch mkt.Type {
	case core.MarketTypeAsk:
		bid, err := s.catalogStg.Index(mkt.ItemID)
		if err != nil {
			return err
		}
		if bid.Quantity != 0 && bid.HighestBid > mkt.Price {
			return core.MarketErrInvalidAskPrice
		}
	case core.MarketTypeBid:
		ask, err := s.catalogStg.Index(mkt.ItemID)
		if err != nil {
			return err
		}
		if ask.Quantity != 0 && ask.LowestAsk < mkt.Price {
			return core.MarketErrInvalidBidPrice
		}
	}

	return nil
}

func (s *marketService) Update(ctx context.Context, mkt *core.Market) error {
	cur, err := s.checkOwnership(ctx, mkt.ID)
	if err != nil {
		return err
	}
	// Prevents access to update existing market when account is flagged.
	if err = s.checkFlaggedUser(cur.UserID); err != nil {
		return err
	}

	if err = mkt.CheckUpdate(); err != nil {
		return err
	}

	// Resolves steam profile URL input as partner steam id.
	if strings.TrimSpace(mkt.PartnerSteamID) != "" {
		mkt.PartnerSteamID, err = s.steam.ResolveVanityURL(mkt.PartnerSteamID)
		if err != nil {
			return err
		}
	}
	// Try to find a matching bid and set its status to complete.
	if cur.Type == core.MarketTypeAsk && mkt.Status == core.MarketStatusReserved {
		if err = s.AutoCompleteBid(ctx, *cur, mkt.PartnerSteamID); err != nil {
			return err
		}
	}

	// Append note to existing notes.
	mkt.Notes = strings.TrimSpace(fmt.Sprintf("%s\n%s", cur.Notes, mkt.Notes))

	// Do not allow update on these fields.
	mkt.UserID = ""
	mkt.ItemID = ""
	mkt.Price = 0
	mkt.Currency = ""
	if err = s.marketStg.Update(mkt); err != nil {
		return err
	}

	go func() {
		if _, err = s.catalogStg.Index(mkt.ItemID); err != nil {
			s.logger.Errorf("could not index item %s: %s", mkt.ItemID, err)
		}
	}()

	if mkt.Type == core.MarketTypeAsk {
		switch mkt.Status {
		case core.MarketStatusReserved:
			s.dispatch.VerifyInventory(mkt.ID)
		case core.MarketStatusSold:
			s.dispatch.VerifyDelivery(mkt.ID)
		}
	}

	s.getRelatedFields(mkt)
	return nil
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

// AutoCompleteBid detects if there's matching reservation on buy order and automatically
// resolve it by setting complete-bid status.
func (s *marketService) AutoCompleteBid(_ context.Context, ask core.Market, partnerSteamID string) error {
	if ask.ItemID == "" || ask.UserID == "" || partnerSteamID == "" {
		return fmt.Errorf("ask market item id, user id, and partner steam id are required")
	}

	// Use buyer ID to get the matching market.
	buyer, err := s.userStg.Get(partnerSteamID)
	if err != nil {
		return nil
	}

	// Find matching bid market to update status.
	fo := core.FindOpts{
		Filter: core.Market{
			Type:   core.MarketTypeBid,
			Status: core.MarketStatusLive,
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
	b.Status = core.MarketStatusBidCompleted
	b.PartnerSteamID = seller.SteamID
	return s.marketStg.Update(&b)
}

func (s *marketService) checkOwnership(ctx context.Context, id string) (*core.Market, error) {
	au := core.AuthFromContext(ctx)
	if au == nil {
		return nil, core.AuthErrNoAccess
	}

	mkt, err := s.userMarket(au.UserID, id)
	if err != nil {
		return nil, errors.New(core.AuthErrForbidden, err)
	}

	if mkt == nil {
		return nil, errors.New(core.AuthErrForbidden, core.MarketErrNotFound)
	}

	return mkt, nil
}

func (s *marketService) userMarket(userID, id string) (*core.Market, error) {
	cur, err := s.marketStg.Get(id)
	if err != nil {
		return nil, err
	}
	if cur.UserID != userID {
		return nil, core.MarketErrNotFound
	}

	return cur, nil
}

func (s *marketService) Catalog(opts core.FindOpts) ([]core.Catalog, *core.FindMetadata, error) {
	res, err := s.catalogStg.Find(opts)
	if err != nil {
		return nil, nil, err
	}

	if !opts.WithMeta {
		return res, nil, err
	}

	// Get result and total count for metadata.
	tc, err := s.catalogStg.Count(opts)
	if err != nil {
		return nil, nil, err
	}

	return res, &core.FindMetadata{
		ResultCount: len(res),
		TotalCount:  tc,
	}, nil
}

func (s *marketService) TrendingCatalog(opts core.FindOpts) ([]core.Catalog, *core.FindMetadata, error) {
	res, err := s.catalogStg.Trending()
	if err != nil {
		return nil, nil, err
	}

	if !opts.WithMeta {
		return res, nil, err
	}

	return res, &core.FindMetadata{
		ResultCount: len(res),
		TotalCount:  10, // Fixed value of top 10
	}, nil
}

func (s *marketService) CatalogDetails(slug string) (*core.Catalog, error) {
	if slug == "" {
		return nil, core.CatalogErrNotFound
	}

	c, err := s.catalogStg.Get(slug)
	if err == core.CatalogErrNotFound {
		i, err := s.itemStg.GetBySlug(slug)
		if err != nil {
			return nil, err
		}

		c := i.ToCatalog()
		return &c, nil

	} else if err != nil {
		return nil, err
	}

	// Retrieve 10 live asks entries.
	mf := core.Market{Type: core.MarketTypeAsk, ItemID: c.ID, Status: core.MarketStatusLive}
	fo := core.FindOpts{
		Filter: mf,
		Limit:  10,
		Sort:   "price",
	}
	res, _, err := s.Markets(context.Background(), fo)
	if err != nil {
		return nil, err
	}
	c.Asks = res
	// Retrieve 10 live bids entries.
	//mf.Type = core.MarketTypeBid
	//fo.Filter = mf
	//fo.Sort = "price"
	//fo.Desc = true
	//res, _, err = s.Markets(context.Background(), fo)
	//if err != nil {
	//	return nil, err
	//}
	//c.Bids = res

	return c, err
}
