package jobs

import (
	"context"
	"time"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/gokit/log"
	"github.com/kudarap/dotagiftx/steaminv"
	"github.com/kudarap/dotagiftx/verified"
)

// VerifyInventory represents a inventory verification job.
type VerifyInventory struct {
	inventorySvc core.InventoryService
	marketSvc    core.MarketService
	logger       log.Logger

	interval time.Duration
}

func NewVerifyInventory(is core.InventoryService, ms core.MarketService, lg log.Logger) *VerifyInventory {
	return &VerifyInventory{is, ms, lg, defaultJobInterval}
}

func (i *VerifyInventory) String() string { return "verify_inventory" }

func (i *VerifyInventory) Interval() time.Duration { return i.interval }

func (i *VerifyInventory) Run(ctx context.Context) error {
	opts := core.FindOpts{Filter: core.Market{Type: core.MarketTypeAsk, Status: core.MarketStatusLive}}
	opts.Sort = "updated_at:desc"
	opts.Limit = 10
	opts.Page = 1

	src := steaminv.InventoryAsset
	for {
		res, _, err := i.marketSvc.Markets(ctx, opts)
		if err != nil {
			return err
		}

		for _, mkt := range res {
			// Skip verified statuses.
			if mkt.InventoryStatus == core.InventoryStatusVerified ||
				mkt.InventoryStatus == core.InventoryStatusNoHit {
				continue
			}

			status, assets, err := verified.Inventory(src, mkt.User.SteamID, mkt.Item.Name)
			if err != nil {
				continue
			}
			i.logger.Println("batch", opts.Page, mkt.User.SteamID, mkt.Item.Name, status)

			err = i.inventorySvc.Set(ctx, &core.Inventory{
				MarketID: mkt.ID,
				Status:   status,
				Assets:   assets,
			})
			if err != nil {
				i.logger.Errorln(mkt.User.SteamID, mkt.Item.Name, status, err)
			}
		}

		// should continue batching?
		if len(res) == 0 {
			return nil
		}
		opts.Page++
	}
}
