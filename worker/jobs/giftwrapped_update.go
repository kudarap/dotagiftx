package jobs

import (
	"context"
	"log/slog"
	"time"

	"github.com/kudarap/dotagiftx/dotagiftx"
	"github.com/kudarap/dotagiftx/verify"
)

// GiftWrappedUpdate represents a job that will update delivered items that still unopened.
type GiftWrappedUpdate struct {
	deliverySvc  deliveryService
	deliveryRepo deliveryRepository
	marketRepo   marketRepository
	source       *verify.Source
	logger       *slog.Logger
	// job settings
	name     string
	interval time.Duration
	filter   dotagiftx.Delivery
}

func NewGiftWrappedUpdate(
	ds deliveryService,
	dg deliveryRepository,
	ms marketRepository,
	vs *verify.Source,
	lg *slog.Logger,
) *GiftWrappedUpdate {
	falsePtr := false
	f := dotagiftx.Delivery{
		GiftOpened: &falsePtr,
		Status:     dotagiftx.DeliveryStatusSenderVerified,
	}
	return &GiftWrappedUpdate{
		ds, dg, ms, vs, lg,
		"giftwrapped_update", time.Hour / 2, f}
}

func (gw *GiftWrappedUpdate) String() string { return gw.name }

func (gw *GiftWrappedUpdate) Interval() time.Duration { return gw.interval }

func (gw *GiftWrappedUpdate) Run(ctx context.Context) error {
	bs := time.Now()
	defer func() {
		gw.logger.Info("GIFT WRAPPED UPDATE BENCHMARK TIME", "elapsed", time.Since(bs))
	}()

	opts := dotagiftx.FindOpts{Filter: gw.filter}
	opts.Sort = "updated_at:desc"
	opts.Limit = 10
	opts.Page = 0
	opts.IndexKey = "status"

	for {
		deliveries, err := gw.deliveryRepo.ToVerify(ctx, opts)
		if err != nil {
			return err
		}

		for _, dd := range deliveries {
			start := time.Now()

			gw.logger.Info("processing gift wrapped update", "delivery_id", dd.ID, "gift_opened", *dd.GiftOpened, "retries", dd.Retries)
			if dd.RetriesExceeded() {
				continue
			}

			mkt, _ := gw.market(ctx, dd.MarketID)
			if mkt == nil {
				gw.logger.Error("skipped process! market not found")
				continue
			}

			if mkt.User == nil || mkt.Item == nil {
				gw.logger.Error("skipped process! missing data", "user", mkt.User, "item", mkt.Item)
				continue
			}

			result, err := gw.source.Delivery(ctx, mkt.User.Name, mkt.PartnerSteamID, mkt.Item.Name)
			if err != nil {
				gw.logger.Error("delivery verification error", "error", err)
				continue
			}
			gw.logger.Info("batch", "page", opts.Page, "user", mkt.User.Name, "partner_steam_id", mkt.PartnerSteamID, "item", mkt.Item.Name, "status", result.Status)

			err = gw.deliverySvc.Set(ctx, &dotagiftx.Delivery{
				MarketID:   mkt.ID,
				Status:     result.Status,
				Assets:     result.Assets,
				VerifiedBy: result.VerifiedBy,
				ElapsedMs:  time.Since(start).Milliseconds(),
			})
			if err != nil {
				gw.logger.Error("delivery update error", "steam_id", mkt.User.SteamID, "item", mkt.Item.Name, "status", result.Status, "error", err)
			}

			// rest(5)
		}

		// Is there more?
		if len(deliveries) < opts.Limit {
			return nil
		}
		// opts.Page++
	}
}

func (gw *GiftWrappedUpdate) market(ctx context.Context, id string) (*dotagiftx.Market, error) {
	f := dotagiftx.FindOpts{Filter: dotagiftx.Market{ID: id}}
	markets, err := gw.marketRepo.Find(ctx, f)
	if err != nil {
		return nil, err
	}
	if len(markets) == 0 {
		return nil, nil
	}
	mkt := markets[0]
	return &mkt, nil
}
