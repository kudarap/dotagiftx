package jobs

import (
	"context"
	"time"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/logging"
	"github.com/kudarap/dotagiftx/verify"
)

// GiftWrappedUpdate represents a job that will update delivered items that still unopened.
type GiftWrappedUpdate struct {
	deliverySvc dotagiftx.DeliveryService
	deliveryStg dotagiftx.DeliveryStorage
	marketStg   dotagiftx.MarketStorage
	source      *verify.Source
	logger      logging.Logger
	// job settings
	name     string
	interval time.Duration
	filter   dotagiftx.Delivery
}

func NewGiftWrappedUpdate(
	ds dotagiftx.DeliveryService,
	dg dotagiftx.DeliveryStorage,
	ms dotagiftx.MarketStorage,
	vs *verify.Source,
	lg logging.Logger,
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
		gw.logger.Println("GIFT WRAPPED UPDATE BENCHMARK TIME", time.Since(bs))
	}()

	opts := dotagiftx.FindOpts{Filter: gw.filter}
	opts.Sort = "updated_at:desc"
	opts.Limit = 10
	opts.Page = 0
	opts.IndexKey = "status"

	for {
		deliveries, err := gw.deliveryStg.ToVerify(opts)
		if err != nil {
			return err
		}

		for _, dd := range deliveries {
			start := time.Now()

			gw.logger.Infoln("processing gift wrapped update", dd.ID, *dd.GiftOpened, dd.Retries)
			if dd.RetriesExceeded() {
				continue
			}

			mkt, _ := gw.market(dd.MarketID)
			if mkt == nil {
				gw.logger.Errorf("skipped process! market not found")
				continue
			}

			if mkt.User == nil || mkt.Item == nil {
				gw.logger.Errorf("skipped process! missing data user:%#v item:%#v", mkt.User, mkt.Item)
				continue
			}

			result, err := gw.source.Delivery(ctx, mkt.User.Name, mkt.PartnerSteamID, mkt.Item.Name)
			if err != nil {
				gw.logger.Errorf("delivery verification error: %s", err)
				continue
			}
			gw.logger.Println("batch", opts.Page, mkt.User.Name, mkt.PartnerSteamID, mkt.Item.Name, result.Status)

			err = gw.deliverySvc.Set(ctx, &dotagiftx.Delivery{
				MarketID:   mkt.ID,
				Status:     result.Status,
				Assets:     result.Assets,
				VerifiedBy: result.VerifiedBy,
				ElapsedMs:  time.Since(start).Milliseconds(),
			})
			if err != nil {
				gw.logger.Errorln(mkt.User.SteamID, mkt.Item.Name, result.Status, err)
			}

			//rest(5)
		}

		// Is there more?
		if len(deliveries) < opts.Limit {
			return nil
		}
		//opts.Page++
	}
}

func (gw *GiftWrappedUpdate) market(id string) (*dotagiftx.Market, error) {
	f := dotagiftx.FindOpts{Filter: dotagiftx.Market{ID: id}}
	markets, err := gw.marketStg.Find(f)
	if err != nil {
		return nil, err
	}
	if len(markets) == 0 {
		return nil, nil
	}
	mkt := markets[0]
	return &mkt, nil
}
