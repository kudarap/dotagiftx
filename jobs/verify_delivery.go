package jobs

import (
	"context"
	"time"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/gokit/log"
	"github.com/kudarap/dotagiftx/steaminv"
	"github.com/kudarap/dotagiftx/verified"
)

// VerifyDelivery represents a delivery verification job.
type VerifyDelivery struct {
	marketSvc   core.MarketService
	deliverySvc core.DeliveryService
	logger      log.Logger
}

func NewVerifyDelivery(ms core.MarketService, ds core.DeliveryService, lg log.Logger) *VerifyDelivery {
	return &VerifyDelivery{ms, ds, lg}
}

func (j *VerifyDelivery) String() string { return "verify_delivery" }

func (j *VerifyDelivery) Interval() time.Duration { return time.Hour }

func (j *VerifyDelivery) Run(ctx context.Context) error {
	opts := core.FindOpts{Filter: core.Market{Type: core.MarketTypeAsk, Status: core.MarketStatusSold}}
	opts.Sort = "updated_at:desc"
	opts.Limit = 10
	opts.Page = 1

	src := steaminv.InventoryAsset
	for {
		res, _, err := j.marketSvc.Markets(ctx, opts)
		if err != nil {
			return err
		}

		for _, mkt := range res {
			// Skip verified statuses.
			if mkt.DeliveryStatus == core.DeliveryStatusNameVerified ||
				mkt.DeliveryStatus == core.DeliveryStatusSenderVerified {
				continue
			}

			status, items, err := verified.Delivery(src, mkt.User.Name, mkt.PartnerSteamID, mkt.Item.Name)
			if err != nil {
				continue
			}

			j.logger.Println("batch", opts.Page, mkt.User.SteamID, mkt.Item.Name, status, len(items))
		}

		// should continue batching?
		if len(res) == 0 {
			return nil
		}
		opts.Page++
	}
}
