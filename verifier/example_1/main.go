package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/steam"
	"github.com/kudarap/dotagiftx/steaminv"
)

func main() {
	//inv, err := steaminv.SWR("76561198088587178")
	//fmt.Println(inv, err)

	delivered, _ := getDelivered()
	verifiedDelivery(delivered)
}

func verifiedDelivery(markets []core.Market) {
	var processed, failed, verified int

	// Benchmark things up.
	ts := time.Now()
	defer func() {
		fmt.Println(time.Now().Sub(ts))
	}()

	for _, mkt := range markets {
		processed++
		fmt.Println(strings.Repeat("-", 70))
		fmt.Println(fmt.Sprintf("%s -> %s (%s)", mkt.User.Name, mkt.PartnerSteamID, mkt.Item.Name))
		fmt.Println(strings.Repeat("-", 70))

		res, err := verify(mkt.User.Name, mkt.PartnerSteamID, mkt.Item.Name)
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Println("")
			failed++
			continue
		}

		fmt.Println("Found:", len(res))
		if len(res) != 0 {
			r := res[0]
			fmt.Println("GiftFrom:", r.GiftFrom)
			fmt.Println("DateReceived:", r.DateReceived)
			fmt.Println("Dedication:", r.Dedication)
			verified++
		}

		fmt.Println("")
	}

	fmt.Println(fmt.Sprintf("%d/%d total | %d error | %d/%d verified", processed, len(markets), failed, processed-verified, verified))
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
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	return data.Data, nil
}

func verify(sellerPersona, buyerSteamID, itemName string) ([]steam.Asset, error) {
	inv, err := steaminv.SWR(buyerSteamID)
	if err != nil {
		return nil, fmt.Errorf("could not get inventory: %s", err)
	}
	if inv == nil {
		return nil, fmt.Errorf("inventory empty result")
	}

	var fi []steam.Asset
	for _, inv := range inv.ToAssets() {
		// Checking against seller persona name might not be accurate since
		// buyer can clear gift information that's why it need to snapshot buyer
		// inventory immediately.
		if inv.GiftFrom != sellerPersona {
			//continue
		}

		// Checks target item name from description and name field.
		if !strings.Contains(strings.Join(inv.Descriptions, "|"), itemName) &&
			!strings.Contains(inv.Name, itemName) {
			continue
		}

		fi = append(fi, inv)
	}

	return fi, nil
}
