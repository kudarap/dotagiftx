package service

import (
	"context"
	"errors"

	dgx "github.com/kudarap/dotagiftx"
)

var ErrHammerNotWeilded = errors.New("user is not weilding a hmmer")

// markedOfBaal special number to detect eternal mark of doom.
const markedOfBaal = 10000

// NewHammerService returns a new Ban service.
func NewHammerService(us dgx.UserStorage, ms dgx.MarketStorage) *BanService {
	return &BanService{us, ms}
}

type BanService struct {
	userStg   dgx.UserStorage
	marketStg dgx.MarketStorage
}

func (s *BanService) Ban(ctx context.Context, p dgx.HammerParams) (*dgx.User, error) {
	return s.hilt(ctx, p, dgx.UserStatusBanned)
}

func (s *BanService) Suspend(ctx context.Context, p dgx.HammerParams) (*dgx.User, error) {
	return s.hilt(ctx, p, dgx.UserStatusSuspended)
}

func (s *BanService) Lift(ctx context.Context, steamID string, restoreListings bool) error {
	au := dgx.AuthFromContext(ctx)
	if au == nil {
		return dgx.AuthErrNoAccess
	}
	if err := s.weildingHammer(au.UserID); err != nil {
		return err
	}

	u, err := s.userStg.Get(steamID)
	if err != nil {
		return err
	}
	u.Status += markedOfBaal // Marked! I could use this to track what was the last offense.
	if err := s.userStg.Update(u); err != nil {
		return err
	}

	if !restoreListings {
		return nil
	}

	// Listing restoration
	return s.restoreListings(u.ID)
}

func (s *BanService) hilt(ctx context.Context, p dgx.HammerParams, us dgx.UserStatus) (*dgx.User, error) {
	au := dgx.AuthFromContext(ctx)
	if au == nil {
		return nil, dgx.AuthErrNoAccess
	}
	if err := s.weildingHammer(au.UserID); err != nil {
		return nil, err
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}

	u, err := s.userStg.Get(p.SteamID)
	if err != nil {
		return nil, err
	}

	u.Status = us
	u.Notes = p.Reason
	if err := s.userStg.Update(u); err != nil {
		return nil, err
	}

	if err := s.cancelListings(u.ID); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *BanService) cancelListings(userID string) error {
	return s.sunderListings(userID, dgx.MarketStatusLive, dgx.MarketStatusCancelled)
}

func (s *BanService) restoreListings(userID string) error {
	return s.sunderListings(userID, dgx.MarketStatusCancelled, dgx.MarketStatusLive)
}

func (s *BanService) sunderListings(userID string, from, to dgx.MarketStatus) error {
	f := dgx.Market{
		UserID: userID,
		Status: from,
	}
	ms, err := s.marketStg.Find(dgx.FindOpts{Filter: f})
	if err != nil {
		return err
	}

	for _, mm := range ms {
		mm.Status = to
		if err := s.marketStg.BaseUpdate(&mm); err != nil {
			return err
		}
	}
	return nil
}

func (s *BanService) weildingHammer(userID string) error {
	u, err := s.userStg.Get(userID)
	if err != nil {
		return err
	}

	if u == nil || !u.Hammer {
		return ErrHammerNotWeilded
	}
	return nil
}
