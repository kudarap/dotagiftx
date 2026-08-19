package dotagiftx

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// markOfBaal special number to detect eternal mark of doom.
const markOfBaal = 10000

var ErrHammerNotWielded = errors.New("user is not wielding a hammer")

// HammerParams represents parameters to drop some suspension and bans.
type HammerParams struct {
	SteamID string `json:"steam_id"`
	Reason  string `json:"reason"`
}

func (p HammerParams) Validate() error {
	if strings.TrimSpace(p.SteamID) == "" && strings.TrimSpace(p.Reason) == "" {
		return fmt.Errorf("steamd_id and reason is required")
	}

	return nil
}

// NewHammerService returns a new Ban service.
func NewHammerService(us userRepository, ms marketRepository) *BanService {
	return &BanService{us, ms}
}

type BanService struct {
	userRepo   userRepository
	marketRepo marketRepository
}

func (s *BanService) Ban(ctx context.Context, p HammerParams) (*User, error) {
	return s.hilt(ctx, p, UserStatusBanned)
}

func (s *BanService) Suspend(ctx context.Context, p HammerParams) (*User, error) {
	return s.hilt(ctx, p, UserStatusSuspended)
}

func (s *BanService) Lift(ctx context.Context, steamID string, restoreListings bool) error {
	au := AuthFromContext(ctx)
	if au == nil {
		return AuthErrNoAccess
	}
	if err := s.wieldingHammer(ctx, au.UserID); err != nil {
		return err
	}

	u, err := s.userRepo.Get(ctx, steamID)
	if err != nil {
		return err
	}
	u.Status += markOfBaal // Marked! I could use this to track what was the last offense.
	if err := s.userRepo.Update(ctx, u); err != nil {
		return err
	}

	if !restoreListings {
		return nil
	}

	// Listing restoration
	return s.restoreListings(ctx, u.ID)
}

func (s *BanService) hilt(ctx context.Context, p HammerParams, us UserStatus) (*User, error) {
	au := AuthFromContext(ctx)
	if au == nil {
		return nil, AuthErrNoAccess
	}
	if err := s.wieldingHammer(ctx, au.UserID); err != nil {
		return nil, err
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}

	u, err := s.userRepo.Get(ctx, p.SteamID)
	if err != nil {
		return nil, err
	}

	u.Status = us
	u.Notes = p.Reason
	if err := s.userRepo.Update(ctx, u); err != nil {
		return nil, err
	}

	if err := s.cancelListings(ctx, u.ID); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *BanService) cancelListings(ctx context.Context, userID string) error {
	return s.sunderListings(ctx, userID, MarketStatusLive, MarketStatusCancelled)
}

func (s *BanService) restoreListings(ctx context.Context, userID string) error {
	return s.sunderListings(ctx, userID, MarketStatusCancelled, MarketStatusLive)
}

func (s *BanService) sunderListings(ctx context.Context, userID string, from, to MarketStatus) error {
	f := Market{
		UserID: userID,
		Status: from,
	}
	ms, err := s.marketRepo.Find(ctx, FindOpts{Filter: f})
	if err != nil {
		return err
	}

	for _, mm := range ms {
		mm.Status = to
		if err := s.marketRepo.BaseUpdate(ctx, &mm); err != nil {
			return err
		}
	}
	return nil
}

func (s *BanService) wieldingHammer(ctx context.Context, userID string) error {
	u, err := s.userRepo.Get(ctx, userID)
	if err != nil {
		return err
	}

	if u == nil || !u.Hammer {
		return ErrHammerNotWielded
	}
	return nil
}
