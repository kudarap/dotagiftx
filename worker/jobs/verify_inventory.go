package jobs

import (
	"context"
	"time"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/logging"
	"github.com/kudarap/dotagiftx/phantasm"
	"github.com/kudarap/dotagiftx/steaminvorg"
	"github.com/kudarap/dotagiftx/verify"
)

// VerifyInventory represents an inventory verification job.
type VerifyInventory struct {
	inventorySvc dotagiftx.InventoryService
	marketStg    dotagiftx.MarketStorage
	phantasmSvc  *phantasm.Service
	logger       logging.Logger
	// job settings
	name     string
	interval time.Duration
	filter   dotagiftx.Market
}

func NewVerifyInventory(is dotagiftx.InventoryService, ms dotagiftx.MarketStorage, ps *phantasm.Service, lg logging.Logger) *VerifyInventory {
	f := dotagiftx.Market{}
	return &VerifyInventory{
		is, ms, ps, lg,
		"verify_inventory", time.Hour * 24, f}
}

func (vi *VerifyInventory) String() string { return vi.name }

func (vi *VerifyInventory) Interval() time.Duration { return vi.interval }

func (vi *VerifyInventory) Run(ctx context.Context) error {
	bs := time.Now()
	defer func() {
		vi.logger.Println("VERIFIED INVENTORY BENCHMARK TIME", time.Since(bs))
	}()

	opts := dotagiftx.FindOpts{Filter: vi.filter}
	opts.IndexSorting = true
	opts.Sort = "updated_at"
	opts.Desc = true
	opts.Limit = 10
	opts.Page = 0

	source := steaminvorg.InventoryAssetWithCache
	source = vi.phantasmSvc.InventoryAsset
	for {
		res, err := vi.marketStg.PendingInventoryStatus(opts)
		if err != nil {
			return err
		}

		for _, mkt := range res {
			// Skip verified statuses.
			if mkt.InventoryStatus == dotagiftx.InventoryStatusVerified ||
				mkt.InventoryStatus == dotagiftx.InventoryStatusNoHit {

				// TODO! might remove items
				//vi.logger.Warnln("batch no need check", opts.Page, mkt.User.SteamID, mkt.Item.Name)
				continue
			}

			if mkt.User == nil || mkt.Item == nil {
				vi.logger.Errorf("skipped process! missing data user:%#v item:%#v", mkt.User, mkt.Item)
				continue
			}

			status, assets, err := verify.Inventory(source, mkt.User.SteamID, mkt.Item.Name)
			if err != nil {
				continue
			}
			vi.logger.Println("batch", opts.Page, mkt.User.SteamID, mkt.Item.Name, status)

			err = vi.inventorySvc.Set(ctx, &dotagiftx.Inventory{
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
