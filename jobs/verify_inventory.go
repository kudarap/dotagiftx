package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/gokit/log"
	"github.com/kudarap/dotagiftx/steaminv"
	"github.com/kudarap/dotagiftx/verified"
)

// VerifyInventory represents a inventory verification job.
type VerifyInventory struct {
	inventorySvc core.InventoryService
	marketStg    core.MarketStorage
	logger       log.Logger
	// job settings
	name     string
	interval time.Duration
	filter   core.Market
}

func NewVerifyInventory(is core.InventoryService, ms core.MarketStorage, lg log.Logger) *VerifyInventory {
	f := core.Market{Type: core.MarketTypeAsk, Status: core.MarketStatusLive}
	return &VerifyInventory{
		is, ms, lg,
		"verify_inventory", defaultJobInterval, f}
}

func (vi *VerifyInventory) String() string { return vi.name }

func (vi *VerifyInventory) Interval() time.Duration { return vi.interval }

func (vi *VerifyInventory) Run(ctx context.Context) error {
	bs := time.Now()
	defer func() {
		fmt.Println("======== VERIFIED INVENTORY BENCHMARK TIME =========")
		fmt.Println(time.Now().Sub(bs))
		fmt.Println("====================================================")
	}()

	opts := core.FindOpts{Filter: vi.filter}
	opts.Sort = "updated_at:desc"
	opts.Limit = 10
	opts.Page = 0

	src := steaminv.InventoryAssetWithCache
	for {
		res, err := vi.marketStg.PendingInventoryStatus(opts)
		if err != nil {
			return err
		}

		for _, mkt := range res {
			// Skip verified statuses.
			if mkt.InventoryStatus == core.InventoryStatusVerified ||
				mkt.InventoryStatus == core.InventoryStatusNoHit {

				// TODO! might remove items
				vi.logger.Warnln("batch no need check", opts.Page, mkt.User.SteamID, mkt.Item.Name)
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

			//rest(5)
			time.Sleep(time.Second / 4)
		}

		// Is there more?
		if len(res) < opts.Limit {
			return nil
		}
		//opts.Page++
		//time.Sleep(time.Second * 2)
	}
}
