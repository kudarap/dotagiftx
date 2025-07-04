package jobs

import (
	"context"
	"time"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/logging"
	"github.com/kudarap/dotagiftx/verify"
)

// VerifyDelivery represents a delivery verification job.
type VerifyDelivery struct {
	deliverySvc dotagiftx.DeliveryService
	marketStg   dotagiftx.MarketStorage
	source      *verify.Source
	logger      logging.Logger
	// job settings
	name     string
	interval time.Duration
	filter   dotagiftx.Market
}

func NewVerifyDelivery(
	ds dotagiftx.DeliveryService,
	ms dotagiftx.MarketStorage,
	vs *verify.Source,
	lg logging.Logger,
) *VerifyDelivery {
	f := dotagiftx.Market{Type: dotagiftx.MarketTypeAsk, Status: dotagiftx.MarketStatusSold}
	return &VerifyDelivery{
		ds, ms, vs, lg,
		"verify_delivery", time.Hour * 12, f}
}

func (vd *VerifyDelivery) String() string { return vd.name }

func (vd *VerifyDelivery) Interval() time.Duration { return vd.interval }

func (vd *VerifyDelivery) Run(ctx context.Context) error {
	bs := time.Now()
	defer func() {
		vd.logger.Println("VERIFIED DELIVERY BENCHMARK TIME", time.Since(bs))
	}()

	opts := dotagiftx.FindOpts{Filter: vd.filter}
	opts.IndexSorting = true
	opts.Sort = "updated_at"
	opts.Desc = true
	opts.Limit = 10
	opts.Page = 0

	for {
		res, err := vd.marketStg.PendingDeliveryStatus(opts)
		if err != nil {
			return err
		}

		for _, mkt := range res {
			start := time.Now()

			// Skip verified statuses.
			if mkt.DeliveryStatus == dotagiftx.DeliveryStatusNameVerified ||
				mkt.DeliveryStatus == dotagiftx.DeliveryStatusSenderVerified {
				continue
			}

			if mkt.User == nil || mkt.Item == nil {
				vd.logger.Errorf("skipped process! missing data user:%#v item:%#v", mkt.User, mkt.Item)
				continue
			}

			result, err := vd.source.Delivery(ctx, mkt.User.Name, mkt.PartnerSteamID, mkt.Item.Name)
			if err != nil {
				continue
			}

			vd.logger.Println("batch", opts.Page, mkt.User.Name, mkt.PartnerSteamID, mkt.Item.Name, result.Status)
			err = vd.deliverySvc.Set(ctx, &dotagiftx.Delivery{
				MarketID:   mkt.ID,
				Status:     result.Status,
				Assets:     result.Assets,
				VerifiedBy: result.VerifiedBy,
				ElapsedMs:  time.Since(start).Milliseconds(),
			})
			if err != nil {
				vd.logger.Errorln(mkt.User.SteamID, mkt.Item.Name, result.Status, err)
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
