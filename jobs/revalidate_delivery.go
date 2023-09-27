package jobs

import (
	"context"
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
		"revalidate_delivery", time.Hour, f}
}

func (rd *RevalidateDelivery) String() string { return rd.name }

func (rd *RevalidateDelivery) Interval() time.Duration { return rd.interval }

func (rd *RevalidateDelivery) Run(ctx context.Context) error {
	bs := time.Now()
	defer func() {
		rd.logger.Println("REVALIDATE DELIVERY BENCHMARK TIME", time.Since(bs))
	}()

	opts := core.FindOpts{Filter: rd.filter}
	opts.Sort = "updated_at:desc"
	opts.Limit = 10
	opts.Page = 0
	opts.IndexKey = "status"

	src := steaminv.InventoryAsset
	for {
		res, err := rd.marketStg.PendingDeliveryStatus(opts)
		if err != nil {
			return err
		}

		for _, mkt := range res {
			if mkt.User == nil || mkt.Item == nil {
				rd.logger.Debug("skipped process! missing data user:%#v item:%#v", mkt.User, mkt.Item)
				continue
			}

			if mkt.Delivery == nil {
				rd.logger.Debug("skipped process! no delivery data")
				continue
			}
			if mkt.Delivery.Retries > 10 {
				rd.logger.Debug("skipped process! max retries reached")
				continue
			}

			status, assets, err := verified.Delivery(src, mkt.User.Name, mkt.PartnerSteamID, mkt.Item.Name)
			if err != nil {
				continue
			}
			rd.logger.Println("batch", opts.Page, mkt.User.Name, mkt.PartnerSteamID, mkt.Item.Name, status)

			err = rd.deliverySvc.Set(ctx, &core.Delivery{
				MarketID: mkt.ID,
				Status:   status,
				Assets:   assets,
			})
			if err != nil {
				rd.logger.Errorln(mkt.User.SteamID, mkt.Item.Name, status, err)
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
