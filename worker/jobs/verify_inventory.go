package jobs

import (
	"context"
	"time"

	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/gokit/log"
	"github.com/kudarap/dotagiftx/steaminvorg"
	"github.com/kudarap/dotagiftx/verifying"
)

// VerifyInventory represents inventory verification job.
type VerifyInventory struct {
	inventorySvc dgx.InventoryService
	marketStg    dgx.MarketStorage
	logger       log.Logger
	// job settings
	name     string
	interval time.Duration
	filter   dgx.Market
}

func NewVerifyInventory(is dgx.InventoryService, ms dgx.MarketStorage, lg log.Logger) *VerifyInventory {
	f := dgx.Market{}
	return &VerifyInventory{
		is, ms, lg,
		"verify_inventory", time.Hour * 24, f}
}

func (vi *VerifyInventory) String() string { return vi.name }

func (vi *VerifyInventory) Interval() time.Duration { return vi.interval }

func (vi *VerifyInventory) Run(ctx context.Context) error {
	bs := time.Now()
	defer func() {
		vi.logger.Println("VERIFIED INVENTORY BENCHMARK TIME", time.Since(bs))
	}()

	opts := dgx.FindOpts{Filter: vi.filter}
	opts.Sort = "updated_at:desc"
	opts.Limit = 10
	opts.Page = 0

	source := steaminvorg.InventoryAssetWithCache
	for {
		res, err := vi.marketStg.PendingInventoryStatus(opts)
		if err != nil {
			return err
		}

		for _, mkt := range res {
			// Skip verified statuses.
			if mkt.InventoryStatus == dgx.InventoryStatusVerified ||
				mkt.InventoryStatus == dgx.InventoryStatusNoHit {

				// TODO! might remove items
				//vi.logger.Warnln("batch no need check", opts.Page, mkt.User.SteamID, mkt.Item.Name)
				continue
			}

			if mkt.User == nil || mkt.Item == nil {
				vi.logger.Errorf("skipped process! missing data user:%#v item:%#v", mkt.User, mkt.Item)
				continue
			}

			status, assets, err := verifying.Inventory(source, mkt.User.SteamID, mkt.Item.Name)
			if err != nil {
				continue
			}
			vi.logger.Println("batch", opts.Page, mkt.User.SteamID, mkt.Item.Name, status)

			err = vi.inventorySvc.Set(ctx, &dgx.Inventory{
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
