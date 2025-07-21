package legacy

import (
	"context"
	"errors"

	"github.com/kudarap/dotagiftx"
)

var ErrHammerNotWielded = errors.New("user is not wielding a hammer")

// markedOfBaal special number to detect eternal mark of doom.
const markedOfBaal = 10000

// NewHammerService returns a new Ban service.
func NewHammerService(us dotagiftx.UserStorage, ms dotagiftx.MarketStorage) *BanService {
	return &BanService{us, ms}
}

type BanService struct {
	userStg   dotagiftx.UserStorage
	marketStg dotagiftx.MarketStorage
}

func (s *BanService) Ban(ctx context.Context, p dotagiftx.HammerParams) (*dotagiftx.User, error) {
	return s.hilt(ctx, p, dotagiftx.UserStatusBanned)
}

func (s *BanService) Suspend(ctx context.Context, p dotagiftx.HammerParams) (*dotagiftx.User, error) {
	return s.hilt(ctx, p, dotagiftx.UserStatusSuspended)
}

func (s *BanService) Lift(ctx context.Context, steamID string, restoreListings bool) error {
	au := dotagiftx.AuthFromContext(ctx)
	if au == nil {
		return dotagiftx.AuthErrNoAccess
	}
	if err := s.wieldingHammer(au.UserID); err != nil {
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

func (s *BanService) hilt(ctx context.Context, p dotagiftx.HammerParams, us dotagiftx.UserStatus) (*dotagiftx.User, error) {
	au := dotagiftx.AuthFromContext(ctx)
	if au == nil {
		return nil, dotagiftx.AuthErrNoAccess
	}
	if err := s.wieldingHammer(au.UserID); err != nil {
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
	return s.sunderListings(userID, dotagiftx.MarketStatusLive, dotagiftx.MarketStatusCancelled)
}

func (s *BanService) restoreListings(userID string) error {
	return s.sunderListings(userID, dotagiftx.MarketStatusCancelled, dotagiftx.MarketStatusLive)
}

func (s *BanService) sunderListings(userID string, from, to dotagiftx.MarketStatus) error {
	f := dotagiftx.Market{
		UserID: userID,
		Status: from,
	}
	ms, err := s.marketStg.Find(dotagiftx.FindOpts{Filter: f})
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

func (s *BanService) wieldingHammer(userID string) error {
	u, err := s.userStg.Get(userID)
	if err != nil {
		return err
	}

	if u == nil || !u.Hammer {
		return ErrHammerNotWielded
	}
	return nil
}
