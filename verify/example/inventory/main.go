package main

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/kudarap/dotagiftx/config"
	"github.com/kudarap/dotagiftx/phantasm"
	"github.com/kudarap/dotagiftx/redis"
	"github.com/kudarap/dotagiftx/steaminvorg"
	"github.com/kudarap/dotagiftx/verify"
)

func main() {
	var conf config.Config
	if err := config.Load(&conf); err != nil {
		panic("could not load config: " + err.Error())
	}
	redisClient, err := redis.New(conf.Redis)
	if err != nil {
		panic(err)
	}

	logger := slog.Default()
	slog.SetLogLoggerLevel(slog.LevelDebug)

	phantasmSvc := phantasm.NewService(conf.Phantasm, redisClient, logger)
	assetSrc := verify.JoinAssetSource(
		phantasmSvc.InventoryAssetWithProvider,
		steaminvorg.InventoryAssetWithProvider,
	)

	params := []struct {
		steamID, item string
	}{
		{"76561198088587178", "Dirge Amplifier"},
		{"76561198088587178", "Fluttering Breeze"},
		{"76561198078663607", "Loaded Prospects"},
	}

	ctx := context.Background()
	for _, param := range params {
		result, err := verify.Inventory(ctx, assetSrc, param.steamID, param.item)
		fmt.Println(strings.Repeat("-", 70))
		fmt.Printf("%s -> %s\n", param.steamID, param.item)
		fmt.Println(strings.Repeat("-", 70))
		if err != nil {
			fmt.Printf("Errored: %s \n\n", err)
			continue
		}

		snaps := result.Assets
		fmt.Println("Verified by:", result.VerifiedBy)
		fmt.Println("Status:", result.Status)
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
