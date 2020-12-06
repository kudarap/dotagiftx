package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/errors"
	"github.com/sirupsen/logrus"
)

// NewMarket returns new Market service.
func NewMarket(
	ss core.MarketStorage,
	us core.UserStorage,
	is core.ItemStorage,
	ts core.TrackStorage,
	cs core.CatalogStorage,
	sc core.SteamClient,
	lg *logrus.Logger,
) core.MarketService {
	return &marketService{ss, us, is, ts, cs, sc, lg}
}

type marketService struct {
	marketStg  core.MarketStorage
	userStg    core.UserStorage
	itemStg    core.ItemStorage
	trackStg   core.TrackStorage
	catalogStg core.CatalogStorage
	steam      core.SteamClient
	logger     *logrus.Logger
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
	return mkt, nil
}

func (s *marketService) getRelatedFields(mkt *core.Market) {
	mkt.User, _ = s.userStg.Get(mkt.UserID)
	mkt.Item, _ = s.itemStg.Get(mkt.ItemID)
}

func (s *marketService) Create(ctx context.Context, mkt *core.Market) error {
	// Set market ownership.
	au := core.AuthFromContext(ctx)
	if au == nil {
		return core.AuthErrNoAccess
	}
	mkt.UserID = au.UserID

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

	return nil
}

func (s *marketService) checkAskType(ask *core.Market) error {
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
	return nil
}

func (s *marketService) Update(ctx context.Context, mkt *core.Market) error {
	cur, err := s.checkOwnership(ctx, mkt.ID)
	if err != nil {
		return err
	}

	if err := mkt.CheckUpdate(); err != nil {
		return err
	}

	// Resolved steam profile URL input as partner steam id.
	if mkt.Status == core.MarketStatusReserved {
		mkt.PartnerSteamID, err = s.steam.ResolveVanityURL(mkt.PartnerSteamID)
		if err != nil {
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
	if err := s.marketStg.Update(mkt); err != nil {
		return err
	}

	go func() {
		if _, err := s.catalogStg.Index(mkt.ItemID); err != nil {
			s.logger.Errorf("could not index item %s: %s", mkt.ItemID, err)
		}
	}()

	s.getRelatedFields(mkt)
	return nil
}

func (s *marketService) checkOwnership(ctx context.Context, id string) (*core.Market, error) {
	au := core.AuthFromContext(ctx)
	if au == nil {
		return nil, core.AuthErrNoAccess
	}

	mkt, err := s.userMarket(au.UserID, id)
	if err != nil {
		return nil, errors.New(core.AuthErrNoAccess, err)
	}

	if mkt == nil {
		return nil, errors.New(core.AuthErrNoAccess, core.MarketErrNotFound)
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

	return c, err
}
