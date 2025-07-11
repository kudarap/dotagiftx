package jobs

import (
	"context"
	"time"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/logging"
	"github.com/kudarap/dotagiftx/verify"
)

// RevalidateDelivery represents a delivery verification job.
type RevalidateDelivery struct {
	deliverySvc dotagiftx.DeliveryService
	marketStg   dotagiftx.MarketStorage
	source      *verify.Source
	logger      logging.Logger
	// job settings
	name     string
	interval time.Duration
	filter   dotagiftx.Market
}

func NewRevalidateDelivery(
	ds dotagiftx.DeliveryService,
	ms dotagiftx.MarketStorage,
	vs *verify.Source,
	lg logging.Logger,
) *RevalidateDelivery {
	f := dotagiftx.Market{Type: dotagiftx.MarketTypeAsk, Status: dotagiftx.MarketStatusSold}
	return &RevalidateDelivery{
		ds, ms, vs, lg,
		"revalidate_delivery", time.Hour * 12, f}
}

func (rd *RevalidateDelivery) String() string { return rd.name }

func (rd *RevalidateDelivery) Interval() time.Duration { return rd.interval }

func (rd *RevalidateDelivery) Run(ctx context.Context) error {
	bs := time.Now()
	defer func() {
		rd.logger.Println("REVALIDATE DELIVERY BENCHMARK TIME", time.Since(bs))
	}()

	opts := dotagiftx.FindOpts{Filter: rd.filter}
	opts.Sort = "updated_at:desc"
	opts.Limit = 10
	opts.Page = 0
	opts.IndexKey = "status"

	for {
		res, err := rd.marketStg.PendingDeliveryStatus(opts)
		if err != nil {
			return err
		}

		for _, mkt := range res {
			start := time.Now()

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

			result, err := rd.source.Delivery(ctx, mkt.User.Name, mkt.PartnerSteamID, mkt.Item.Name)
			if err != nil {
				continue
			}
			rd.logger.Println("batch", opts.Page, mkt.User.Name, mkt.PartnerSteamID, mkt.Item.Name, result.Status)

			err = rd.deliverySvc.Set(ctx, &dotagiftx.Delivery{
				MarketID:   mkt.ID,
				Status:     result.Status,
				Assets:     result.Assets,
				VerifiedBy: result.VerifiedBy,
				ElapsedMs:  time.Since(start).Milliseconds(),
			})
			if err != nil {
				rd.logger.Errorln(mkt.User.SteamID, mkt.Item.Name, result.Status, err)
			}

			//rest(5)
			//time.Sleep(time.Second / 4)
		}

		// Is there more?
		if len(res) < opts.Limit {
			return nil
		}
		//opts.Page++
	}
}
