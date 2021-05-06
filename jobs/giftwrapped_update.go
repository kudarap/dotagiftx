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

// GiftWrappedUpdate represents a job that will update delivered
// items that still un-opened
type GiftWrappedUpdate struct {
	deliverySvc core.DeliveryService
	deliveryStg core.DeliveryStorage
	marketStg   core.MarketStorage
	logger      log.Logger
	// job settings
	name     string
	interval time.Duration
	filter   core.Delivery
}

func NewGiftWrappedUpdate(ds core.DeliveryService, dstg core.DeliveryStorage, ms core.MarketStorage, lg log.Logger) *GiftWrappedUpdate {
	falsePtr := false
	f := core.Delivery{
		GiftOpened: &falsePtr,
		Status:     core.DeliveryStatusSenderVerified,
	}
	return &GiftWrappedUpdate{
		ds, dstg, ms, lg,
		"giftwrapped_update", time.Hour, f}
}

func (vd *GiftWrappedUpdate) String() string { return vd.name }

func (vd *GiftWrappedUpdate) Interval() time.Duration { return vd.interval }

func (vd *GiftWrappedUpdate) Run(ctx context.Context) error {
	bs := time.Now()
	defer func() {
		fmt.Println("======== GIFT WRAPPED UPDATE BENCHMARK TIME =========")
		fmt.Println(time.Now().Sub(bs))
		fmt.Println("====================================================")
	}()

	opts := core.FindOpts{Filter: vd.filter}
	opts.Sort = "updated_at:desc"
	opts.Limit = 10
	opts.Page = 0

	src := steaminv.InventoryAsset
	for {
		deliveries, err := vd.deliveryStg.ToVerify(opts)
		if err != nil {
			return err
		}

		for _, dd := range deliveries {
			vd.logger.Infoln("processing gift wrapped update", dd.ID, *dd.GiftOpened, dd.Retries)
			if dd.RetriesExceeded() {
				continue
			}

			mkt, _ := vd.market(dd.MarketID)
			if mkt == nil {
				vd.logger.Errorf("skipped process! market not found")
				continue
			}

			if mkt.User == nil || mkt.Item == nil {
				vd.logger.Errorf("skipped process! missing data user:%#v item:%#v", mkt.User, mkt.Item)
				continue
			}

			status, assets, err := verified.Delivery(src, mkt.User.Name, mkt.PartnerSteamID, mkt.Item.Name)
			if err != nil {
				vd.logger.Errorf("delivery verification error: %s", err)
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

			//rest(5)
		}

		// Is there more?
		if len(deliveries) < opts.Limit {
			return nil
		}
		//opts.Page++
	}
}

func (vd *GiftWrappedUpdate) market(id string) (*core.Market, error) {
	f := core.FindOpts{Filter: core.Market{ID: id}}
	markets, err := vd.marketStg.Find(f)
	if err != nil {
		return nil, err
	}
	if len(markets) == 0 {
		return nil, nil
	}
	mkt := markets[0]
	return &mkt, nil
}
