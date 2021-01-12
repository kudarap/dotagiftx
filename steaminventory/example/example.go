package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/kudarap/dotagiftx/steaminventory"

	"github.com/kudarap/dotagiftx/core"
)

func main2() {
	//status, err := steaminventory.Crawl("76561198088587178")
	//fmt.Println(status, err)

	//meta, err := steaminventory.GetMeta("76561198032085715")
	//fmt.Println(meta, err)

	//inv, err := steaminventory.Get("76561198088587178")
	//fmt.Println(inv, err)

	//inv, err := steaminventory.GetNWait("76561198088587178")
	//fmt.Println(inv, err)

	//flat, err := steaminventory.NewFlatInventoryFromV2(*inv)
	//fmt.Println(flat, err)

}

func main() {
	//flat, err := steaminventory.VerifyDelivery("karosu!", "76561198088587178", "Ravenous Abyss")
	//fmt.Println(flat, err)

	delivered, _ := getDelivered()
	for _, mkt := range delivered {
		res, err := steaminventory.VerifyDelivery(mkt.User.Name, mkt.PartnerSteamID, mkt.Item.Name)
		if err != nil {
			fmt.Println("could not verify!", err)
			continue
		}

		fmt.Println(strings.Repeat("-", 70))
		fmt.Println(fmt.Sprintf("%s -> %s (%s)", mkt.User.Name, mkt.PartnerSteamID, mkt.Item.Name))
		fmt.Println(strings.Repeat("-", 70))
		fmt.Println("Found:", len(res))
		if len(res) != 0 {
			r := res[0]
			fmt.Println("GiftFrom:", r.GiftFrom)
			fmt.Println("DateReceived:", r.DateReceived)
			fmt.Println("Dedication:", r.Dedication)
		}

		fmt.Println("")
	}
}

func getDelivered() ([]core.Market, error) {
	resp, err := http.Get("https://api.dotagiftx.com/markets?sort=updated_at:desc&limit=50&status=400")
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
