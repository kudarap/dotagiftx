package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/kudarap/dotagiftx/steaminv"
	"github.com/kudarap/dotagiftx/verified"
)

func main() {
	assetSrc := steaminv.InventoryAsset

	seller := "tuty # aka FUTUX"
	buyerID := "76561198265102770"
	item := "Fissured Flight"

	status, snaps, err := verified.Delivery(assetSrc, seller, buyerID, item)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Status:", status)
	fmt.Println("Items:", len(snaps))

	for _, as := range snaps {
		fmt.Println(strings.Repeat("-", 70))
		fmt.Println("GiftFrom:", as.GiftFrom)
		fmt.Println("DateReceived:", as.DateReceived)
		fmt.Println("Dedication:", as.Dedication)
	}
}
