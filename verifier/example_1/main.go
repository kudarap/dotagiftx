package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/steaminv"
	"github.com/kudarap/dotagiftx/verifier"
)

func main() {
	assetSrc := steaminv.InventoryAsset

	var errorCtr, okCtr, privateCtr, noHitCtr, itemCtr, sellerCtr int

	// Benchmark things up.
	ts := time.Now()
	defer func() {
		fmt.Println(time.Now().Sub(ts))
	}()

	items, _ := getDelivered()

	for _, item := range items {
		status, snaps, err := verifier.Delivery(assetSrc, item.User.Name, item.PartnerSteamID, item.Item.Name)

		fmt.Println(strings.Repeat("-", 70))
		fmt.Println(fmt.Sprintf("%s -> %s (%s)", item.User.Name, item.PartnerSteamID, item.Item.Name))
		fmt.Println(strings.Repeat("-", 70))
		fmt.Println("Status:", status)
		if err != nil {
			errorCtr++
			fmt.Printf("Errored: %s \n\n", err)
			return
			//continue
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
		case verifier.VerifyStatusPrivate:
			privateCtr++
		case verifier.VerifyStatusNoHit:
			noHitCtr++
		case verifier.VerifyStatusItem:
			itemCtr++
		case verifier.VerifyStatusSeller:
			sellerCtr++
		}

		fmt.Println("")
	}

	fmt.Println(fmt.Sprintf("%d/%d total | %d error", okCtr, len(items), errorCtr))
	fmt.Println(fmt.Sprintf("%d private | %d nohit | %d item | %d seller",
		privateCtr, noHitCtr, itemCtr, sellerCtr))

}

func getDelivered() ([]core.Market, error) {
	resp, err := http.Get("https://api.dotagiftx.com/markets?sort=updated_at:desc&limit=1000&status=400")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data := struct {
		Data []core.Market
	}{}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	return data.Data, nil
}
