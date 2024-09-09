package jobs

import (
	"context"
	"time"

	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/gokit/log"
	"github.com/kudarap/dotagiftx/steaminvorg"
	"github.com/kudarap/dotagiftx/verifying"
)

// RevalidateDelivery represents a delivery verification job.
type RevalidateDelivery struct {
	deliverySvc dgx.DeliveryService
	marketStg   dgx.MarketStorage
	logger      log.Logger
	// job settings
	name     string
	interval time.Duration
	filter   dgx.Market
}

func NewRevalidateDelivery(ds dgx.DeliveryService, ms dgx.MarketStorage, lg log.Logger) *RevalidateDelivery {
	f := dgx.Market{Type: dgx.MarketTypeAsk, Status: dgx.MarketStatusSold}
	return &RevalidateDelivery{
		ds, ms, lg,
		"revalidate_delivery", time.Hour * 24, f}
}

func (rd *RevalidateDelivery) String() string { return rd.name }

func (rd *RevalidateDelivery) Interval() time.Duration { return rd.interval }

func (rd *RevalidateDelivery) Run(ctx context.Context) error {
	bs := time.Now()
	defer func() {
		rd.logger.Println("REVALIDATE DELIVERY BENCHMARK TIME", time.Since(bs))
	}()

	opts := dgx.FindOpts{Filter: rd.filter}
	opts.Sort = "updated_at:desc"
	opts.Limit = 10
	opts.Page = 0
	opts.IndexKey = "status"

	src := steaminvorg.InventoryAsset
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

			status, assets, err := verifying.Delivery(src, mkt.User.Name, mkt.PartnerSteamID, mkt.Item.Name)
			if err != nil {
				continue
			}
			rd.logger.Println("batch", opts.Page, mkt.User.Name, mkt.PartnerSteamID, mkt.Item.Name, status)

			err = rd.deliverySvc.Set(ctx, &dgx.Delivery{
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
