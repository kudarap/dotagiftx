package jobs

import (
	"context"
	"time"

	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/gokit/log"
	"github.com/kudarap/dotagiftx/steaminvorg"
	"github.com/kudarap/dotagiftx/verifying"
)

// VerifyDelivery represents a delivery verification job.
type VerifyDelivery struct {
	deliverySvc dgx.DeliveryService
	marketStg   dgx.MarketStorage
	logger      log.Logger
	// job settings
	name     string
	interval time.Duration
	filter   dgx.Market
}

func NewVerifyDelivery(ds dgx.DeliveryService, ms dgx.MarketStorage, lg log.Logger) *VerifyDelivery {
	f := dgx.Market{Type: dgx.MarketTypeAsk, Status: dgx.MarketStatusSold}
	return &VerifyDelivery{
		ds, ms, lg,
		"verify_delivery", time.Hour * 24, f}
}

func (vd *VerifyDelivery) String() string { return vd.name }

func (vd *VerifyDelivery) Interval() time.Duration { return vd.interval }

func (vd *VerifyDelivery) Run(ctx context.Context) error {
	bs := time.Now()
	defer func() {
		vd.logger.Println("VERIFIED DELIVERY BENCHMARK TIME", time.Since(bs))
	}()

	opts := dgx.FindOpts{Filter: vd.filter}
	opts.Sort = "updated_at:desc"
	opts.Limit = 10
	opts.Page = 0

	src := steaminvorg.InventoryAsset
	for {
		res, err := vd.marketStg.PendingDeliveryStatus(opts)
		if err != nil {
			return err
		}

		for _, mkt := range res {
			// Skip verified statuses.
			if mkt.DeliveryStatus == dgx.DeliveryStatusNameVerified ||
				mkt.DeliveryStatus == dgx.DeliveryStatusSenderVerified {
				continue
			}

			if mkt.User == nil || mkt.Item == nil {
				vd.logger.Errorf("skipped process! missing data user:%#v item:%#v", mkt.User, mkt.Item)
				continue
			}

			status, assets, err := verifying.Delivery(src, mkt.User.Name, mkt.PartnerSteamID, mkt.Item.Name)
			if err != nil {
				continue
			}
			vd.logger.Println("batch", opts.Page, mkt.User.Name, mkt.PartnerSteamID, mkt.Item.Name, status)

			err = vd.deliverySvc.Set(ctx, &dgx.Delivery{
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
