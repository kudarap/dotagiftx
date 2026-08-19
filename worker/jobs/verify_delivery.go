package jobs

import (
	"context"
	"log/slog"
	"time"

	"github.com/kudarap/dotagiftx/dotagiftx"
	"github.com/kudarap/dotagiftx/verify"
)

// VerifyDelivery represents a delivery verification job.
type VerifyDelivery struct {
	deliverySvc deliveryService
	marketRepo  marketRepository
	source      *verify.Source
	logger      *slog.Logger
	// job settings
	name     string
	interval time.Duration
	filter   dotagiftx.Market
}

func NewVerifyDelivery(
	ds deliveryService,
	ms marketRepository,
	vs *verify.Source,
	lg *slog.Logger,
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
		vd.logger.Info("VERIFIED DELIVERY BENCHMARK TIME", "elapsed", time.Since(bs))
	}()

	opts := dotagiftx.FindOpts{Filter: vd.filter}
	opts.IndexSorting = true
	opts.Sort = "updated_at"
	opts.Desc = true
	opts.Limit = 10
	opts.Page = 0

	for {
		res, err := vd.marketRepo.PendingDeliveryStatus(ctx, opts)
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
				vd.logger.Error("skipped process! missing data", "user", mkt.User, "item", mkt.Item)
				continue
			}

			result, err := vd.source.Delivery(ctx, mkt.User.Name, mkt.PartnerSteamID, mkt.Item.Name)
			if err != nil {
				continue
			}

			vd.logger.Info("batch", "page", opts.Page, "user", mkt.User.Name, "partner_steam_id", mkt.PartnerSteamID, "item", mkt.Item.Name, "status", result.Status)
			err = vd.deliverySvc.Set(ctx, &dotagiftx.Delivery{
				MarketID:   mkt.ID,
				Status:     result.Status,
				Assets:     result.Assets,
				VerifiedBy: result.VerifiedBy,
				ElapsedMs:  time.Since(start).Milliseconds(),
			})
			if err != nil {
				vd.logger.Error("delivery update error", "steam_id", mkt.User.SteamID, "item", mkt.Item.Name, "status", result.Status, "error", err)
			}

			time.Sleep(time.Second / 4)
		}

		// Is there more?
		if len(res) < opts.Limit {
			return nil
		}
		// opts.Page++
	}
}
