package runonce

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/kudarap/dotagiftx"
)

// GenerateCompletedBid collects delivered items from partner sellers
// and generate a bid-complted entry from there.
//
// ./fixgenbids 76561198236673500 - for dryrun
// ./fixgenbids -commit 76561198236673500 - commit changes
func GenerateCompletedBid(marketStore dotagiftx.MarketStorage, userStore dotagiftx.UserStorage, marketSvc dotagiftx.MarketService) {
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

	f := dotagiftx.Market{
		Type:           dotagiftx.MarketTypeAsk,
		PartnerSteamID: buyer.SteamID,
	}
	res, err := marketStore.Find(dotagiftx.FindOpts{Filter: f})
	if err != nil {
		log.Println("err", err)
		return
	}

	for _, ask := range res {
		log.Println("===============================================")
		log.Println("generating market bid")

		if ask.Status != dotagiftx.MarketStatusReserved && ask.Status != dotagiftx.MarketStatusSold {
			log.Println("skipping status of", ask.Status)
			continue
		}

		seller, err := userStore.Get(ask.UserID)
		if err != nil {
			log.Fatalln("could not get seller:", err)
		}

		bid := new(dotagiftx.Market)
		bid.UserID = buyer.ID
		bid.Type = dotagiftx.MarketTypeBid
		bid.PartnerSteamID = seller.SteamID
		bid.Status = dotagiftx.MarketStatusBidCompleted
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
		ctx := dotagiftx.AuthToContext(context.TODO(), &dotagiftx.Auth{UserID: bid.UserID})
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
func AutoCompleteBid(marketSvc dotagiftx.MarketService) {
	ctx := context.Background()
	f := dotagiftx.Market{
		Type:   dotagiftx.MarketTypeAsk,
		Status: dotagiftx.MarketStatusReserved,
	}
	res, _, err := marketSvc.Markets(ctx, dotagiftx.FindOpts{Filter: f})
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

func ResolveCompletedBidSteamID(store dotagiftx.MarketStorage, steam dotagiftx.SteamClient) {
	o := dotagiftx.FindOpts{Filter: dotagiftx.Market{Status: dotagiftx.MarketStatusBidCompleted}}
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
			log.Println("could not resolve URL")
			continue
		}

		if err = store.Update(&m); err != nil {
			log.Println("could not update market entry", err)
		}

		log.Println(m.PartnerSteamID, "fixed")
	}
}
