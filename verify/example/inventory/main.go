package main

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kudarap/dotagiftx/phantasm"
	"github.com/kudarap/dotagiftx/verify"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("could not load config: " + err.Error())
	}

	var conf phantasm.Config
	phantasmSvc := phantasm.NewService(conf, slog.Default())

	assetSrc := verify.MultiAssetSource(
		phantasmSvc.InventoryAsset,
		//steaminvorg.InventoryAssetWithCache,
	)

	params := []struct {
		steamID, item string
	}{
		{"76561198088587178", "Dirge Amplifier"},
		{"76561198088587178", "Cannonroar Confessor"},
		{"76561198130214012", "Draca Mane"},
	}

	for _, param := range params {
		status, snaps, err := verify.Inventory(assetSrc, param.steamID, param.item)
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
			fmt.Println("")
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
