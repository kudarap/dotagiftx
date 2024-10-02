package jobs

import (
	"context"
	"time"

	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/gokit/log"
	"github.com/kudarap/dotagiftx/steaminvorg"
	"github.com/kudarap/dotagiftx/verifying"
)

// GiftWrappedUpdate represents a job that will update delivered
// items that still un-opened
type GiftWrappedUpdate struct {
	deliverySvc dgx.DeliveryService
	deliveryStg dgx.DeliveryStorage
	marketStg   dgx.MarketStorage
	logger      log.Logger
	// job settings
	name     string
	interval time.Duration
	filter   dgx.Delivery
}

func NewGiftWrappedUpdate(ds dgx.DeliveryService, dstg dgx.DeliveryStorage, ms dgx.MarketStorage, lg log.Logger) *GiftWrappedUpdate {
	falsePtr := false
	f := dgx.Delivery{
		GiftOpened: &falsePtr,
		Status:     dgx.DeliveryStatusSenderVerified,
	}
	return &GiftWrappedUpdate{
		ds, dstg, ms, lg,
		"giftwrapped_update", time.Hour / 2, f}
}

func (gw *GiftWrappedUpdate) String() string { return gw.name }

func (gw *GiftWrappedUpdate) Interval() time.Duration { return gw.interval }

func (gw *GiftWrappedUpdate) Run(ctx context.Context) error {
	bs := time.Now()
	defer func() {
		gw.logger.Println("GIFT WRAPPED UPDATE BENCHMARK TIME", time.Since(bs))
	}()

	opts := dgx.FindOpts{Filter: gw.filter}
	opts.Sort = "updated_at:desc"
	opts.Limit = 10
	opts.Page = 0
	opts.IndexKey = "status"

	src := steaminvorg.InventoryAsset
	for {
		deliveries, err := gw.deliveryStg.ToVerify(opts)
		if err != nil {
			return err
		}

		for _, dd := range deliveries {
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

			status, assets, err := verifying.Delivery(src, mkt.User.Name, mkt.PartnerSteamID, mkt.Item.Name)
			if err != nil {
				gw.logger.Errorf("delivery verification error: %s", err)
				continue
			}
			gw.logger.Println("batch", opts.Page, mkt.User.Name, mkt.PartnerSteamID, mkt.Item.Name, status)

			err = gw.deliverySvc.Set(ctx, &dgx.Delivery{
				MarketID: mkt.ID,
				Status:   status,
				Assets:   assets,
			})
			if err != nil {
				gw.logger.Errorln(mkt.User.SteamID, mkt.Item.Name, status, err)
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

func (gw *GiftWrappedUpdate) market(id string) (*dgx.Market, error) {
	f := dgx.FindOpts{Filter: dgx.Market{ID: id}}
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
