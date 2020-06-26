package service

import (
	"context"

	"github.com/kudarap/dota2giftables/core"
	"github.com/kudarap/dota2giftables/errors"
)

// NewMarket returns new Market service.
func NewMarket(ss core.MarketStorage, us core.UserStorage, is core.ItemStorage) core.MarketService {
	return &marketService{ss, us, is}
}

type marketService struct {
	marketStg core.MarketStorage
	userStg   core.UserStorage
	itemStg   core.ItemStorage
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
	for i := range res {
		s.getRelatedFields(&res[i])
	}

	if !opts.WithMeta {
		return res, nil, err
	}

	// Get result and total count for metadata.
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
		return errors.New(core.ItemErrRequiredFields, err)
	}

	// Check Item existence.
	if i, _ := s.itemStg.Get(mkt.ItemID); i == nil {
		return core.ItemErrNotFound
	}

	return s.marketStg.Create(mkt)
}

func (s *marketService) Update(ctx context.Context, mkt *core.Market) error {
	_, err := s.checkOwnership(ctx, mkt.ID)
	if err != nil {
		return err
	}

	if err := mkt.CheckUpdate(); err != nil {
		return err
	}

	// Do not allowed update on these fields.
	mkt.UserID = ""
	mkt.ItemID = ""
	mkt.Price = 0
	mkt.Currency = ""
	if err := s.marketStg.Update(mkt); err != nil {
		return err
	}

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

func (s *marketService) Index(opts core.FindOpts) ([]core.MarketIndex, *core.FindMetadata, error) {
	res, err := s.marketStg.FindIndex(opts)
	if err != nil {
		return nil, nil, err
	}

	if !opts.WithMeta {
		return res, nil, err
	}

	// Get result and total count for metadata.
	tc, err := s.marketStg.Count(opts)
	if err != nil {
		return nil, nil, err
	}

	return res, &core.FindMetadata{
		ResultCount: len(res),
		TotalCount:  tc,
	}, nil
}
