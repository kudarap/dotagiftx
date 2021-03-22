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

func (vi *RecheckInventory) String() string { return vi.name }

func (vi *RecheckInventory) Interval() time.Duration { return vi.interval }

func (vi *RecheckInventory) Run(ctx context.Context) error {
	bs := time.Now()
	defer func() {
		fmt.Println("======== RECHECK INVENTORY BENCHMARK TIME =========")
		fmt.Println(time.Now().Sub(bs))
		fmt.Println("====================================================")
	}()

	opts := core.FindOpts{Filter: vi.filter}
	opts.Sort = "updated_at:desc"
	//opts.Limit = 10
	opts.Page = 0

	src := steaminv.InventoryAssetWithCache
	invs, _, err := vi.inventorySvc.Inventories(opts)
	if err != nil {
		return err
	}

	for _, ii := range invs {
		if ii.RetriesExceeded() {
			continue
		}

		mkt, _ := vi.market(ii.MarketID)
		if mkt == nil {
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

	return nil
}

func (vd *RecheckInventory) market(id string) (*core.Market, error) {
	f := core.FindOpts{Filter: core.Market{ID: id}}
	markets, err := vd.marketStg.Find(f)
	if err != nil {
		return nil, err
	}
	if len(markets) == 0 {
		return nil, nil
	}
	mkt := markets[0]
	return &mkt, nil
}
