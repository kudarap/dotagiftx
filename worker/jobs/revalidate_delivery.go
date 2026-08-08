package jobs

import (
	"context"
	"time"

	"log/slog"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/verify"
)

// RevalidateDelivery represents a delivery verification job.
type RevalidateDelivery struct {
	deliverySvc deliveryService
	marketStg   marketStorage
	source      *verify.Source
	logger      *slog.Logger
	// job settings
	name     string
	interval time.Duration
	filter   dotagiftx.Market
}

func NewRevalidateDelivery(
	ds deliveryService,
	ms marketStorage,
	vs *verify.Source,
	lg *slog.Logger,
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
		rd.logger.Info("REVALIDATE DELIVERY BENCHMARK TIME", "elapsed", time.Since(bs))
	}()

	opts := dotagiftx.FindOpts{Filter: rd.filter}
	opts.Sort = "updated_at:desc"
	opts.Limit = 10
	opts.Page = 0
	opts.IndexKey = "status"

	for {
		res, err := rd.marketStg.PendingDeliveryStatus(ctx, opts)
		if err != nil {
			return err
		}

		for _, mkt := range res {
			start := time.Now()

			if mkt.User == nil || mkt.Item == nil {
				rd.logger.Debug("skipped process! missing data", "user", mkt.User, "item", mkt.Item)
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
			rd.logger.Info("batch", "page", opts.Page, "user", mkt.User.Name, "partner_steam_id", mkt.PartnerSteamID, "item", mkt.Item.Name, "status", result.Status)

			err = rd.deliverySvc.Set(ctx, &dotagiftx.Delivery{
				MarketID:   mkt.ID,
				Status:     result.Status,
				Assets:     result.Assets,
				VerifiedBy: result.VerifiedBy,
				ElapsedMs:  time.Since(start).Milliseconds(),
			})
			if err != nil {
				rd.logger.Error("delivery update error", "steam_id", mkt.User.SteamID, "item", mkt.Item.Name, "status", result.Status, "error", err)
			}

			// time.Sleep(time.Second / 4)
		}

		// Is there more?
		if len(res) < opts.Limit {
			return nil
		}
		// opts.Page++
	}
}
