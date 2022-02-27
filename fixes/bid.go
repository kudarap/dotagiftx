package fixes

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kudarap/dotagiftx/core"
)

// GenerateCompletedBid collects delivered items from partner sellers
// and generate a bid-complted entry from there.
//
// ./fixgenbids 76561198236673500 - for dryrun
// ./fixgenbids -commit 76561198236673500 - commit changes
func GenerateCompletedBid(marketStore core.MarketStorage, userStore core.UserStorage, marketSvc core.MarketService) {
	if len(os.Args) < 2 {
		log.Println("steam id argument required")
		os.Exit(1)
	}

	var commit bool
	flag.BoolVar(&commit, "commit", false, "proceed with db changes")
	flag.Parse()

	buyerID := os.Args[len(os.Args)-1]
	log.Println("generating completed bid started", buyerID, "commit:", commit)

	buyer, err := userStore.Get(buyerID)
	if err != nil {
		log.Fatalln("could not get buyer:", err)
	}

	f := core.Market{
		Type:           core.MarketTypeAsk,
		PartnerSteamID: buyer.SteamID,
	}
	res, err := marketStore.Find(core.FindOpts{Filter: f})
	if err != nil {
		log.Println("err", err)
		return
	}

	for _, ask := range res {
		log.Println("===============================================")
		log.Println("generating market bid")

		if ask.Status != core.MarketStatusReserved && ask.Status != core.MarketStatusSold {
			log.Println("skipping status of", ask.Status)
			continue
		}

		seller, err := userStore.Get(ask.UserID)
		if err != nil {
			log.Fatalln("could not get seller:", err)
		}

		bid := new(core.Market)
		bid.UserID = buyer.ID
		bid.Type = core.MarketTypeBid
		bid.PartnerSteamID = seller.SteamID
		bid.Status = core.MarketStatusBidCompleted
		bid.Price = ask.Price
		bid.ItemID = ask.ItemID
		bid.Notes = "test"

		log.Println("ask ID:", ask.ID)
		log.Println("item", ask.Item.Name)
		log.Println("original status:", ask.Status)
		log.Println("-----")
		log.Println("seller ID:", bid.PartnerSteamID)
		log.Println("item ID:", bid.ItemID)
		log.Println("price:", bid.Price)

		if !commit {
			continue
		}
		ctx := core.AuthToContext(context.TODO(), &core.Auth{UserID: bid.UserID})
		if err := marketSvc.Create(ctx, bid); err != nil {
			log.Println("could not create bid:", err)
		}
	}

	log.Println("===============================================")
	log.Println("===============================================")
	if !commit {
		log.Println("NO CHANGES MADE")
	}
	log.Println("generating completed bid finished")
}

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
