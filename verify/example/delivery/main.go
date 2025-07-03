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

	"github.com/joho/godotenv"
	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/phantasm"
	"github.com/kudarap/dotagiftx/steaminvorg"
	"github.com/kudarap/dotagiftx/verify"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("could not load config: " + err.Error())
	}

	var c phantasm.Config
	phantasmSvc := phantasm.NewService(c, slog.Default())

	assetSrc := verify.MultiAssetSource(
		phantasmSvc.InventoryAsset,
		steaminvorg.InventoryAsset,
	)

	var errorCtr, okCtr, privateCtr, noHitCtr, itemCtr, sellerCtr int

	// Benchmark things up.
	ts := time.Now()
	defer func() {
		fmt.Println(time.Now().Sub(ts))
	}()

	items, _ := getDelivered(1)

	ctx := context.Background()
	for _, item := range items {
		status, snaps, err := verify.Delivery(ctx, assetSrc, item.User.Name, item.PartnerSteamID, item.Item.Name)

		fmt.Println(strings.Repeat("-", 70))
		fmt.Println(fmt.Sprintf("%s -> %s (%s)", item.User.Name, item.PartnerSteamID, item.Item.Name))
		fmt.Println(strings.Repeat("-", 70))
		fmt.Println("Status:", status)
		if err != nil {
			errorCtr++
			fmt.Printf("Errored: %s \n\n", err)
			continue
		}

		okCtr++

		fmt.Println("Items:", len(snaps))
		if len(snaps) != 0 {
			r := snaps[0]
			fmt.Println("Name:", r.Name)
			fmt.Println("Contains:", r.Contains)
			fmt.Println("GiftFrom:", r.GiftFrom)
			fmt.Println("DateReceived:", r.DateReceived)
			fmt.Println("Dedication:", r.Dedication)
		}

		switch status {
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

	fmt.Println(fmt.Sprintf("%d/%d total | %d error", okCtr, len(items), errorCtr))
	fmt.Println(fmt.Sprintf("%d private | %d nohit | %d item | %d seller",
		privateCtr, noHitCtr, itemCtr, sellerCtr))

}

func getDelivered(limit int) ([]dotagiftx.Market, error) {
	resp, err := http.Get(fmt.Sprintf(
		"https://api.dotagiftx.com/markets?sort=updated_at:desc&limit=%d&status=400&user_id=%s",
		limit,
		"ddabf335-7286-430a-8403-00e9cda45cfb",
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
