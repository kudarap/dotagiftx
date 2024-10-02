package jobs

import (
	"context"
	"time"

	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/gokit/log"
	"github.com/kudarap/dotagiftx/steaminvorg"
	"github.com/kudarap/dotagiftx/verifying"
)

// RecheckInventory represents a job that rechecks no-hit items.
// Crawling from SteamInventory.org tends to fail sometimes.
type RecheckInventory struct {
	inventorySvc dgx.InventoryService
	marketStg    dgx.MarketStorage
	logger       log.Logger
	// job settings
	name     string
	interval time.Duration
	filter   dgx.Inventory
}

func NewRecheckInventory(is dgx.InventoryService, ms dgx.MarketStorage, lg log.Logger) *RecheckInventory {
	f := dgx.Inventory{Status: dgx.InventoryStatusNoHit}
	return &RecheckInventory{
		is, ms, lg,
		"recheck_inventory", time.Hour, f}
}

func (ri *RecheckInventory) String() string { return ri.name }

func (ri *RecheckInventory) Interval() time.Duration { return ri.interval }

func (ri *RecheckInventory) Run(ctx context.Context) error {
	bs := time.Now()
	defer func() {
		ri.logger.Println("RECHECK INVENTORY BENCHMARK TIME", time.Since(bs))
	}()

	opts := dgx.FindOpts{Filter: ri.filter}
	opts.Sort = "updated_at:desc"
	//opts.Limit = 10
	opts.Page = 0
	opts.IndexKey = "status"

	src := steaminvorg.InventoryAssetWithCache
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

		status, assets, err := verifying.Inventory(src, mkt.User.SteamID, mkt.Item.Name)
		if err != nil {
			continue
		}
		ri.logger.Println("batch", opts.Page, mkt.User.SteamID, mkt.Item.Name, status)

		err = ri.inventorySvc.Set(ctx, &dgx.Inventory{
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

func (ri *RecheckInventory) market(id string) (*dgx.Market, error) {
	f := dgx.FindOpts{Filter: dgx.Market{ID: id}}
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
