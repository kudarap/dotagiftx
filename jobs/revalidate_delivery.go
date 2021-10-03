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

// RevalidateDelivery represents a delivery verification job.
type RevalidateDelivery struct {
	deliverySvc core.DeliveryService
	marketStg   core.MarketStorage
	logger      log.Logger
	// job settings
	name     string
	interval time.Duration
	filter   core.Market
}

func NewRevalidateDelivery(ds core.DeliveryService, ms core.MarketStorage, lg log.Logger) *RevalidateDelivery {
	f := core.Market{Type: core.MarketTypeAsk, Status: core.MarketStatusSold}
	return &RevalidateDelivery{
		ds, ms, lg,
		"verify_delivery", time.Hour, f}
}

func (vd *RevalidateDelivery) String() string { return vd.name }

func (vd *RevalidateDelivery) Interval() time.Duration { return vd.interval }

func (vd *RevalidateDelivery) Run(ctx context.Context) error {
	bs := time.Now()
	defer func() {
		fmt.Println("======== REVALIDATE DELIVERY BENCHMARK TIME =========")
		fmt.Println(time.Since(bs))
		fmt.Println("====================================================")
	}()

	opts := core.FindOpts{Filter: vd.filter}
	opts.Sort = "updated_at:desc"
	opts.Limit = 10
	opts.Page = 0

	src := steaminv.InventoryAsset
	for {
		res, err := vd.marketStg.PendingDeliveryStatus(opts)
		if err != nil {
			return err
		}

		for _, mkt := range res {
			if mkt.User == nil || mkt.Item == nil {
				vd.logger.Errorf("skipped process! missing data user:%#v item:%#v", mkt.User, mkt.Item)
				continue
			}

			status, assets, err := verified.Delivery(src, mkt.User.Name, mkt.PartnerSteamID, mkt.Item.Name)
			if err != nil {
				continue
			}
			vd.logger.Println("batch", opts.Page, mkt.User.Name, mkt.PartnerSteamID, mkt.Item.Name, status)

			err = vd.deliverySvc.Set(ctx, &core.Delivery{
				MarketID: mkt.ID,
				Status:   status,
				Assets:   assets,
			})
			if err != nil {
				vd.logger.Errorln(mkt.User.SteamID, mkt.Item.Name, status, err)
			}

			//rest(5)
			time.Sleep(time.Second / 4)
		}

		// Is there more?
		if len(res) < opts.Limit {
			return nil
		}
		//opts.Page++
	}
}
