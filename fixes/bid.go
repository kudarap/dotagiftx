package fixes

import (
	"context"
	"fmt"
	"log"
	"strings"

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
		if err = marketSvc.AutoCompleteBid(ctx, m, m.PartnerSteamID); err != nil {
			log.Println("could not complete bid", err)
			continue
		}

		log.Println("bid completed", m.ItemID, m.PartnerSteamID)
	}
}

func ResolveCompletedBidSteamID(store core.MarketStorage, steam core.SteamClient) {
	o := core.FindOpts{Filter: core.Market{Status: core.MarketStatusBidCompleted}}
	res, err := store.Find(o)
	if err != nil {
		log.Println("err", err)
		return
	}

	for _, m := range res {
		if !strings.Contains(m.PartnerSteamID, "https://steamcommunity.com") {
			continue
		}

		log.Println("resolving partner steam URL", m.PartnerSteamID)
		m.PartnerSteamID, err = steam.ResolveVanityURL(m.PartnerSteamID)
		if err != nil {
			fmt.Println("could not resolve URL")
			continue
		}

		if err = store.Update(&m); err != nil {
			fmt.Println("could not update market entry", err)
		}

		log.Println(m.PartnerSteamID, "fixed")
	}
}
