package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/steaminvorg"
	"github.com/kudarap/dotagiftx/verifying"
)

func main() {
	assetSrc := steaminvorg.InventoryAsset

	var errorCtr, okCtr, privateCtr, noHitCtr, itemCtr, sellerCtr int

	// Benchmark things up.
	ts := time.Now()
	defer func() {
		fmt.Println(time.Now().Sub(ts))
	}()

	items, _ := getDelivered()

	for _, item := range items {
		status, snaps, err := verifying.Delivery(assetSrc, item.User.Name, item.PartnerSteamID, item.Item.Name)

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
			fmt.Println("GiftFrom:", r.GiftFrom)
			fmt.Println("DateReceived:", r.DateReceived)
			fmt.Println("Dedication:", r.Dedication)
		}

		switch status {
		case dgx.DeliveryStatusPrivate:
			privateCtr++
		case dgx.DeliveryStatusNoHit:
			noHitCtr++
		case dgx.DeliveryStatusNameVerified:
			itemCtr++
		case dgx.DeliveryStatusSenderVerified:
			sellerCtr++
		}

		fmt.Println("")
	}

	fmt.Println(fmt.Sprintf("%d/%d total | %d error", okCtr, len(items), errorCtr))
	fmt.Println(fmt.Sprintf("%d private | %d nohit | %d item | %d seller",
		privateCtr, noHitCtr, itemCtr, sellerCtr))

}

func getDelivered() ([]dgx.Market, error) {
	resp, err := http.Get("https://api.dotagiftx.com/markets?sort=updated_at:desc&limit=1000&status=400")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data := struct {
		Data []dgx.Market
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
