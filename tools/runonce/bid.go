package runonce

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"

	dgx "github.com/kudarap/dotagiftx"
)

// GenerateCompletedBid collects delivered items from partner sellers
// and generate a bid-complted entry from there.
//
// ./fixgenbids 76561198236673500 - for dryrun
// ./fixgenbids -commit 76561198236673500 - commit changes
func GenerateCompletedBid(marketStore dgx.MarketStorage, userStore dgx.UserStorage, marketSvc dgx.MarketService) {
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

	f := dgx.Market{
		Type:           dgx.MarketTypeAsk,
		PartnerSteamID: buyer.SteamID,
	}
	res, err := marketStore.Find(dgx.FindOpts{Filter: f})
	if err != nil {
		log.Println("err", err)
		return
	}

	for _, ask := range res {
		log.Println("===============================================")
		log.Println("generating market bid")

		if ask.Status != dgx.MarketStatusReserved && ask.Status != dgx.MarketStatusSold {
			log.Println("skipping status of", ask.Status)
			continue
		}

		seller, err := userStore.Get(ask.UserID)
		if err != nil {
			log.Fatalln("could not get seller:", err)
		}

		bid := new(dgx.Market)
		bid.UserID = buyer.ID
		bid.Type = dgx.MarketTypeBid
		bid.PartnerSteamID = seller.SteamID
		bid.Status = dgx.MarketStatusBidCompleted
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
		ctx := dgx.AuthToContext(context.TODO(), &dgx.Auth{UserID: bid.UserID})
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
func AutoCompleteBid(marketSvc dgx.MarketService) {
	ctx := context.Background()
	f := dgx.Market{
		Type:   dgx.MarketTypeAsk,
		Status: dgx.MarketStatusReserved,
	}
	res, _, err := marketSvc.Markets(ctx, dgx.FindOpts{Filter: f})
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

func ResolveCompletedBidSteamID(store dgx.MarketStorage, steam dgx.SteamClient) {
	o := dgx.FindOpts{Filter: dgx.Market{Status: dgx.MarketStatusBidCompleted}}
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
