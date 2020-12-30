package fixes

import (
	"context"
	"log"

	"github.com/kudarap/dotagiftx/core"
)

// AutoCompleteBid searches for exiting reservations that has buy order and resolve it.
func AutoCompleteBid(marketSvc core.MarketService) {
	ctx := context.Background()
	f := core.Market{
		Type:   core.MarketTypeAsk,
		Status: core.MarketStatusReserved,
	}
	res, _, err := marketSvc.Markets(ctx, core.FindOpts{Filter: f})
	if err != nil {
		log.Println("err", err)
		return
	}

	for _, m := range res {
		if err = marketSvc.AutoCompleteBid(ctx, m.ItemID, m.PartnerSteamID); err != nil {
			log.Println("could not complete bid", err)
			continue
		}

		log.Println("bid completed", m.ItemID, m.PartnerSteamID)
	}
}
