package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/kudarap/dotagiftx"
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

	phantasmSvc := phantasm.NewService(conf.Phantasm, redisClient, slog.Default())
	assetSrc := verify.JoinAssetSource(
		phantasmSvc.InventoryAssetWithProvider,
		steaminvorg.InventoryAssetWithProvider,
	)

	var errorCtr, okCtr, privateCtr, noHitCtr, itemCtr, sellerCtr int

	// Benchmark things up.
	ts := time.Now()
	defer func() {
		fmt.Println(time.Since(ts))
	}()

	items, _ := getDelivered(1)

	ctx := context.Background()
	for _, item := range items {
		result, err := verify.Delivery(ctx, assetSrc, item.User.Name, item.PartnerSteamID, item.Item.Name)
		fmt.Println(strings.Repeat("-", 70))
		fmt.Printf("%s -> %s (%s)\n", item.User.Name, item.PartnerSteamID, item.Item.Name)
		fmt.Println(strings.Repeat("-", 70))
		if err != nil {
			errorCtr++
			fmt.Printf("Errored: %s \n\n", err)
			continue
		}
		fmt.Println("Verified by:", result.VerifiedBy)
		fmt.Println("Status:", result.Status)

		okCtr++

		snaps := result.Assets
		fmt.Println("Items:", len(snaps))
		if len(snaps) != 0 {
			r := snaps[0]
			fmt.Println("Name:", r.Name)
			fmt.Println("Contains:", r.Contains)
			fmt.Println("GiftFrom:", r.GiftFrom)
			fmt.Println("DateReceived:", r.DateReceived)
			fmt.Println("Dedication:", r.Dedication)
		}

		switch result.Status {
		case dotagiftx.DeliveryStatusPrivate:
			privateCtr++
		case dotagiftx.DeliveryStatusNoHit:
			noHitCtr++
		case dotagiftx.DeliveryStatusNameVerified:
			itemCtr++
		case dotagiftx.DeliveryStatusSenderVerified:
			sellerCtr++
		}

		fmt.Println("")
	}

	fmt.Printf("%d/%d total | %d error\n", okCtr, len(items), errorCtr)
	fmt.Printf("%d private | %d nohit | %d item | %d seller\n", privateCtr, noHitCtr, itemCtr, sellerCtr)
}

func getDelivered(limit int) ([]dotagiftx.Market, error) {
	resp, err := http.Get(fmt.Sprintf(
		"https://api.dotagiftx.com/markets?sort=updated_at:desc&limit=%d&status=400&user_id=%s",
		limit,
		"4a513712-e93e-48fe-bff0-8c653ba50beb",
	))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data := struct {
		Data []dotagiftx.Market
	}{}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	return data.Data, nil
}
