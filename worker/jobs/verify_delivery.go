package jobs

import (
	"context"
	"time"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/logging"
	"github.com/kudarap/dotagiftx/phantasm"
	"github.com/kudarap/dotagiftx/steaminvorg"
	"github.com/kudarap/dotagiftx/verifying"
)

// VerifyDelivery represents a delivery verification job.
type VerifyDelivery struct {
	deliverySvc dotagiftx.DeliveryService
	marketStg   dotagiftx.MarketStorage
	phantasmSvc *phantasm.Service
	logger      logging.Logger
	// job settings
	name     string
	interval time.Duration
	filter   dotagiftx.Market
}

func NewVerifyDelivery(ds dotagiftx.DeliveryService, ms dotagiftx.MarketStorage, ps *phantasm.Service, lg logging.Logger) *VerifyDelivery {
	f := dotagiftx.Market{Type: dotagiftx.MarketTypeAsk, Status: dotagiftx.MarketStatusSold}
	return &VerifyDelivery{
		ds, ms, ps, lg,
		"verify_delivery", time.Hour * 24, f}
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

	src := steaminvorg.InventoryAsset
	src = vd.phantasmSvc.InventoryAsset
	for {
		res, err := vd.marketStg.PendingDeliveryStatus(opts)
		if err != nil {
			return err
		}

		for _, mkt := range res {
			// Skip verified statuses.
			if mkt.DeliveryStatus == dotagiftx.DeliveryStatusNameVerified ||
				mkt.DeliveryStatus == dotagiftx.DeliveryStatusSenderVerified {
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

			err = vd.deliverySvc.Set(ctx, &dotagiftx.Delivery{
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
