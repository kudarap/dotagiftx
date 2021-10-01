package service

import (
	"context"
	"errors"

	"github.com/kudarap/dotagiftx/core"
)

var ErrHammerNotWeilded = errors.New("user is not weilding a hmmer")

// markedOfBaal special number to detect eternal mark of doom.
const markedOfBaal = 10000

// NewHammerService returns a new Ban service.
func NewHammerService(us core.UserStorage, ms core.MarketStorage) *BanService {
	return &BanService{us, ms}
}

type BanService struct {
	userStg   core.UserStorage
	marketStg core.MarketStorage
}

func (s *BanService) Ban(ctx context.Context, p core.HammerParams) (*core.User, error) {
	return s.hilt(ctx, p, core.UserStatusBanned)
}

func (s *BanService) Suspend(ctx context.Context, p core.HammerParams) (*core.User, error) {
	return s.hilt(ctx, p, core.UserStatusSuspended)
}

func (s *BanService) Lift(ctx context.Context, steamID string, restoreListings bool) error {
	au := core.AuthFromContext(ctx)
	if au == nil {
		return core.AuthErrNoAccess
	}

	u, err := s.userStg.Get(steamID)
	if err != nil {
		return err
	}
	if weildingHammer(u) {
		return ErrHammerNotWeilded
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

func (s *BanService) hilt(ctx context.Context, p core.HammerParams, us core.UserStatus) (*core.User, error) {
	au := core.AuthFromContext(ctx)
	if au == nil {
		return nil, core.AuthErrNoAccess
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}

	u, err := s.userStg.Get(p.SteamID)
	if err != nil {
		return nil, err
	}
	if weildingHammer(u) {
		return nil, ErrHammerNotWeilded
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
	return s.sunderListings(userID, core.MarketStatusLive, core.MarketStatusCancelled)
}

func (s *BanService) restoreListings(userID string) error {
	return s.sunderListings(userID, core.MarketStatusCancelled, core.MarketStatusLive)
}

func (s *BanService) sunderListings(userID string, from, to core.MarketStatus) error {
	f := core.Market{
		UserID: userID,
		Status: from,
	}
	ms, err := s.marketStg.Find(core.FindOpts{Filter: f})
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

func weildingHammer(u *core.User) bool {
	if u == nil {
		return false
	}

	return u.Hammer
}
