package jobs

import (
	"context"
	"time"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/gokit/log"
	"github.com/kudarap/dotagiftx/steaminv"
	"github.com/kudarap/dotagiftx/verified"
)

// RecheckInventory represents a job that rechecks no-hit items.
// Crawling from SteamInventory.org tends to fail some times.
type RecheckInventory struct {
	inventorySvc core.InventoryService
	marketStg    core.MarketStorage
	logger       log.Logger
	// job settings
	name     string
	interval time.Duration
	filter   core.Inventory
}

func NewRecheckInventory(is core.InventoryService, ms core.MarketStorage, lg log.Logger) *RecheckInventory {
	f := core.Inventory{Status: core.InventoryStatusNoHit}
	return &RecheckInventory{
		is, ms, lg,
		"recheck_inventory", time.Hour * 12, f}
}

func (ri *RecheckInventory) String() string { return ri.name }

func (ri *RecheckInventory) Interval() time.Duration { return ri.interval }

func (ri *RecheckInventory) Run(ctx context.Context) error {
	bs := time.Now()
	defer func() {
		ri.logger.Println("RECHECK INVENTORY BENCHMARK TIME", time.Since(bs))
	}()

	opts := core.FindOpts{Filter: ri.filter}
	opts.Sort = "updated_at:desc"
	//opts.Limit = 10
	opts.Page = 0

	src := steaminv.InventoryAssetWithCache
	invs, _, err := ri.inventorySvc.Inventories(opts)
	if err != nil {
		return err
	}

	for _, ii := range invs {
		if ii.RetriesExceeded() {
			continue
		}

		mkt, _ := ri.market(ii.MarketID)
		if mkt == nil {
			continue
		}

		if mkt.User == nil || mkt.Item == nil {
			ri.logger.Errorf("skipped process! missing data user:%#v item:%#v", mkt.User, mkt.Item)
			continue
		}

		status, assets, err := verified.Inventory(src, mkt.User.SteamID, mkt.Item.Name)
		if err != nil {
			continue
		}
		ri.logger.Println("batch", opts.Page, mkt.User.SteamID, mkt.Item.Name, status)

		err = ri.inventorySvc.Set(ctx, &core.Inventory{
			MarketID: mkt.ID,
			Status:   status,
			Assets:   assets,
		})
		if err != nil {
			ri.logger.Errorln(mkt.User.SteamID, mkt.Item.Name, status, err)
		}

		//rest(5)
		time.Sleep(time.Second / 4)
	}

	return nil
}

func (ri *RecheckInventory) market(id string) (*core.Market, error) {
	f := core.FindOpts{Filter: core.Market{ID: id}}
	markets, err := ri.marketStg.Find(f)
	if err != nil {
		return nil, err
	}
	if len(markets) == 0 {
		return nil, nil
	}
	mkt := markets[0]
	return &mkt, nil
}
