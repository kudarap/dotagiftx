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
	// job settings
	name     string
	interval time.Duration
	filter   core.Market
}

func NewVerifyInventory(is core.InventoryService, ms core.MarketService, lg log.Logger) *VerifyInventory {
	f := core.Market{Type: core.MarketTypeAsk, Status: core.MarketStatusLive}
	return &VerifyInventory{
		is, ms, lg,
		"verify_inventory", defaultJobInterval, f}
}

func (vi *VerifyInventory) String() string { return vi.name }

func (vi *VerifyInventory) Interval() time.Duration { return vi.interval }

func (vi *VerifyInventory) Run(ctx context.Context) error {
	opts := core.FindOpts{Filter: vi.filter}
	opts.Sort = "updated_at:desc"
	opts.Limit = 10
	opts.Page = 1

	src := steaminv.InventoryAsset
	for {
		res, _, err := vi.marketSvc.Markets(ctx, opts)
		if err != nil {
			return err
		}

		for _, mkt := range res {
			// Skip verified statuses.
			if mkt.InventoryStatus == core.InventoryStatusVerified ||
				mkt.InventoryStatus == core.InventoryStatusNoHit {

				// TODO! might remove items

				continue
			}

			status, assets, err := verified.Inventory(src, mkt.User.SteamID, mkt.Item.Name)
			if err != nil {
				continue
			}
			vi.logger.Println("batch", opts.Page, mkt.User.SteamID, mkt.Item.Name, status)

			err = vi.inventorySvc.Set(ctx, &core.Inventory{
				MarketID: mkt.ID,
				Status:   status,
				Assets:   assets,
			})
			if err != nil {
				vi.logger.Errorln(mkt.User.SteamID, mkt.Item.Name, status, err)
			}
		}

		// should continue batching?
		if len(res) == 0 {
			return nil
		}
		opts.Page++
	}
}
