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

// VerifyDelivery represents a delivery verification job.
type VerifyDelivery struct {
	deliverySvc core.DeliveryService
	marketStg   core.MarketStorage
	logger      log.Logger
	// job settings
	name     string
	interval time.Duration
	filter   core.Market
}

func NewVerifyDelivery(ds core.DeliveryService, ms core.MarketStorage, lg log.Logger) *VerifyDelivery {
	f := core.Market{Type: core.MarketTypeAsk, Status: core.MarketStatusSold}
	return &VerifyDelivery{
		ds, ms, lg,
		"verify_delivery", defaultJobInterval, f}
}

func (vd *VerifyDelivery) String() string { return vd.name }

func (vd *VerifyDelivery) Interval() time.Duration { return vd.interval }

func (vd *VerifyDelivery) Run(ctx context.Context) error {
	bs := time.Now()
	defer func() {
		fmt.Println("======== VERIFIED DELIVERY BENCHMARK TIME =========")
		fmt.Println(time.Now().Sub(bs))
		fmt.Println("====================================================")
	}()

	opts := core.FindOpts{Filter: vd.filter}
	opts.Sort = "updated_at:desc"
	opts.Limit = 10
	opts.Page = 0

	src := steaminv.InventoryAssetWithCache
	for {
		res, err := vd.marketStg.PendingDeliveryStatus(opts)
		if err != nil {
			return err
		}

		for _, mkt := range res {
			// Skip verified statuses.
			if mkt.DeliveryStatus == core.DeliveryStatusNameVerified ||
				mkt.DeliveryStatus == core.DeliveryStatusSenderVerified {
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

			rest(5)
		}

		// Is there more?
		if len(res) < opts.Limit {
			return nil
		}
		//opts.Page++
	}
}
