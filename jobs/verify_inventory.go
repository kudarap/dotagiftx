package jobs

import (
	"context"
	"time"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/gokit/log"
	"github.com/kudarap/dotagiftx/steam"
	"github.com/kudarap/dotagiftx/verified"
)

// VerifyInventory represents a inventory verification job.
type VerifyInventory struct {
	marketSvc core.MarketService
	logger    log.Logger
}

func NewVerifyInventory(ms core.MarketService, lg log.Logger) *VerifyInventory {
	return &VerifyInventory{ms, lg}
}

func (j *VerifyInventory) String() string { return "verify_inventory" }

func (j *VerifyInventory) Interval() time.Duration { return time.Minute * 2 }

func (j *VerifyInventory) Run(ctx context.Context) error {
	o := core.FindOpts{Filter: core.Market{Type: core.MarketTypeAsk, Status: core.MarketStatusLive}}
	res, _, err := j.marketSvc.Markets(ctx, o)
	if err != nil {
		return err
	}

	src := steam.InventoryAsset
	for _, mkt := range res {
		status, items, err := verified.Inventory(src, mkt.User.SteamID, mkt.Item.Name)
		if err != nil {
			continue
		}

		j.logger.Println(mkt.User.SteamID, mkt.Item.Name, status, len(items))
	}

	return nil
}
