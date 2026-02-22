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

// HammerService represents operation for banning and suspending accounts.
type HammerService interface {
	// Ban updates user status to ban and cancels all listings.
	//
	// "Drops the hammer to its eternal doom" is most likely to be permanent.
	Ban(context.Context, HammerParams) (*User, error)

	// Suspend updates user status to suspend and cancels all listings.
	//
	// Fits for those light and abusive offenders. might forget to lift if not reminded.
	Suspend(context.Context, HammerParams) (*User, error)

	// Lift update user status to "marked" and remove its ban or suspend a flag
	// and will restore items if requested.
	Lift(ctx context.Context, steamID string, restoreListings bool) error
}

// NewHammerService returns a new Ban service.
func NewHammerService(us UserStorage, ms MarketStorage) *BanService {
	return &BanService{us, ms}
}

type BanService struct {
	userStg   UserStorage
	marketStg MarketStorage
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
	if err := s.wieldingHammer(au.UserID); err != nil {
		return err
	}

	u, err := s.userStg.Get(steamID)
	if err != nil {
		return err
	}
	u.Status += markOfBaal // Marked! I could use this to track what was the last offense.
	if err := s.userStg.Update(u); err != nil {
		return err
	}

	if !restoreListings {
		return nil
	}

	// Listing restoration
	return s.restoreListings(u.ID)
}

func (s *BanService) hilt(ctx context.Context, p HammerParams, us UserStatus) (*User, error) {
	au := AuthFromContext(ctx)
	if au == nil {
		return nil, AuthErrNoAccess
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
	return s.sunderListings(userID, MarketStatusLive, MarketStatusCancelled)
}

func (s *BanService) restoreListings(userID string) error {
	return s.sunderListings(userID, MarketStatusCancelled, MarketStatusLive)
}

func (s *BanService) sunderListings(userID string, from, to MarketStatus) error {
	f := Market{
		UserID: userID,
		Status: from,
	}
	ms, err := s.marketStg.Find(FindOpts{Filter: f})
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
