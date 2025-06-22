package main

import (
	"fmt"
	"strings"

	"github.com/kudarap/dotagiftx/steaminvorg"
	"github.com/kudarap/dotagiftx/verifying"
)

func main() {
	assetSrc := steaminvorg.InventoryAssetWithCache
	//assetSrc = steam.InventoryAssetWithCache

	params := []struct {
		steamID, item string
	}{
		{"76561198088587178", "Tribal Pathways"},
		{"76561198088587178", "Cannonroar Confessor"},
	}

	for _, param := range params {
		status, snaps, err := verifying.Inventory(assetSrc, param.steamID, param.item)

		fmt.Println(strings.Repeat("-", 70))
		fmt.Println(fmt.Sprintf("%s -> %s", param.steamID, param.item))
		fmt.Println(strings.Repeat("-", 70))
		fmt.Println("Status:", status)
		if err != nil {
			fmt.Printf("Errored: %s \n\n", err)
			continue
		}

		fmt.Println("Items:", len(snaps))
		if len(snaps) == 0 {
			continue
		}

		r := snaps[0]
		fmt.Println("GiftFrom:", r.GiftFrom)
		fmt.Println("DateReceived:", r.DateReceived)
		fmt.Println("Dedication:", r.Dedication)
		for _, ss := range snaps {
			fmt.Println(ss.Name, "qty:", ss.Qty)
		}
		fmt.Println("")
	}
}
