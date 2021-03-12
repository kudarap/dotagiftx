package jobs

import (
	"context"
	"log"
	"time"

	"github.com/kudarap/dotagiftx/steam"
	"github.com/kudarap/dotagiftx/verified"

	"github.com/kudarap/dotagiftx/core"
)

// VerifyInventory represents a inventory verification job.
type VerifyInventory struct {
	marketSvc core.MarketService
}

func NewVerifyInventory(marketSvc core.MarketService) *VerifyInventory {
	return &VerifyInventory{marketSvc}
}

func (v *VerifyInventory) String() string { return "verify_inventory" }

func (v *VerifyInventory) Interval() time.Duration { return time.Minute * 2 }

func (v *VerifyInventory) Run(ctx context.Context) error {
	o := core.FindOpts{Filter: core.Market{Type: core.MarketTypeAsk, Status: core.MarketStatusLive}}
	res, _, err := v.marketSvc.Markets(ctx, o)
	if err != nil {
		return err
	}

	src := steam.InventoryAsset
	for _, mkt := range res {
		status, items, err := verified.Inventory(src, mkt.User.SteamID, mkt.Item.Name)
		if err != nil {
			continue
		}

		log.Println(mkt.User.SteamID, mkt.Item.Name, status, len(items))
	}

	return nil
}
