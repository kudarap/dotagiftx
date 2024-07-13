package dgx

import (
	"context"
	"fmt"
	"strings"
)

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
	// Ban updates user sttus to banned and cancels all listings.
	//
	// "Drops the hammer to its eternal doom" its most likely to be permanent.
	Ban(context.Context, HammerParams) (*User, error)

	// Suspend updates user sttus to suspended and cancels all listings.
	//
	// Fits for those light and abusive offenders. might forget to lift if not reminded.
	Suspend(context.Context, HammerParams) (*User, error)

	// Lift update user status to "marked" and remove its ban or suspend flag
	// and will restore items if requested.
	Lift(ctx context.Context, steamID string, restoreListings bool) error
}
