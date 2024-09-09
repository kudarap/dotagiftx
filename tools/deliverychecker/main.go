package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/kudarap/dotagiftx/steaminvorg"
	"github.com/kudarap/dotagiftx/verifying"
)

func main() {
	assetSrc := steaminvorg.InventoryAsset

	var sellerName string
	var buyerSteamID string
	var itemName string
	flag.StringVar(&sellerName, "s", "", "seller name")
	flag.StringVar(&buyerSteamID, "b", "", "buyer steam id")
	flag.StringVar(&itemName, "i", "", "item name")
	flag.Parse()

	status, snaps, err := verifying.Delivery(assetSrc, sellerName, buyerSteamID, itemName)
	fmt.Println(strings.Repeat("-", 70))
	fmt.Println(fmt.Sprintf("%s -> %s (%s)", sellerName, buyerSteamID, itemName))
	fmt.Println(strings.Repeat("-", 70))
	fmt.Println("Status:", status)
	if err != nil {
		fmt.Printf("Errored: %s \n\n", err)
	}

	fmt.Println("Items:", len(snaps))
	if len(snaps) != 0 {
		r := snaps[0]
		fmt.Println("GiftFrom:", r.GiftFrom)
		fmt.Println("DateReceived:", r.DateReceived)
		fmt.Println("Dedication:", r.Dedication)
	}
}
