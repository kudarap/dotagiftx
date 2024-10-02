package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/steaminvorg"
	"github.com/kudarap/dotagiftx/verifying"
)

func main() {
	assetSrc := steaminvorg.InventoryAssetWithCache

	params := []struct {
		steamID, item string
	}{
		{"76561198011836301", "Dipper the Destroyer"},
		{"76561198011836301", "Dipper the Destroyer Bundle"},

		//{"76561198287849998", "Awaleb's Trundleweed"},
		//{"76561198287849998", "Golden Awaleb's Trundleweed"},

		{"76561198309634734", "The Abscesserator Bundle"},
		{"76561198255419442", "The Abscesserator Bundle"},

		//{"76561198209729119", "Pyrexaec Forge"},
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
			fmt.Println(ss.Name, ss.Qty)
		}

		fmt.Println("")
	}
}

func main2() {
	assetSrc := steaminvorg.InventoryAsset

	params := []struct {
		persona, steamID, item string
	}{
		//{"kudarap", "765611____1477544", "BAD ID"},
		//{"kudarap", "76561198011477544", "PRIVATE"},
		//{"kudarap", "76561198042690669", "NO HIT"},

		{"no luck no win", "76561198046903060", "The Abscesserator Bundle"},

		//{"L3G1ON", "76561198052842962", "Pattern of the Silken Queen"},
		//{"L3G1ON", "76561198052842962", "Evolution of the Infinite"},
		//{"L3G1ON", "76561198052842962", "Pyrexaec Forge"},
		//{"L3G1ON", "76561198052842962", "Secrets of the Celestial"},
		//{"L3G1ON", "76561198052842962", "Shimmer of the Anointed"},

		//{"kudarap", "76561198042690669", "Riddle of the Hierophant"},
		//{"kudarap", "76561198436826874", "Fowl Omen"},
		//{"kudarap", "76561198170463425", "Cinder Sensei"},
		//{"kudarap", "76561198872556187", "Grim Destiny"},
		//{"kudarap", "76561198042690669", "Scorched Amber"},
		//{"kudarap", "76561198872556187", "Tales of the Windward Rogue"},
		//{"kudarap", "76561198203634725", "Adornments of the Jade Emissary"},
		//{"gippeum", "76561198088587178", "Elements of the Endless Plane"},
		//
		//{"Berserk", "76561198355627060", "Shattered Greatsword"},
		//{"Berserk", "76561198042690669", "Ancient Inheritance"},
		//{"Berserk", "76561198042690669", "Poacher's Bane"},
		//{"Berserk", "76561198042690669", "Allure of the Faeshade Flower"},
		//{"Berserk", "76561198256569879", "Endless Night"},
		//{"Berserk", "76561198139657329", "Glimmer of the Sacred Hunt"},
		//
		//{"karosu!", "76561198088587178", "Ravenous Abyss"},
		//{"Araragi-", "76561198809365008", "Master of the Searing Path"},
		//{"Accel", "76561198042690669", "Forsworn Legacy"},
		//{"Dark Knight", "76561198116319576", "Legends of Darkheart Pursuit"},
		//{"ZAAALLGO", "76561198203634725", "Cunning Corsair"},
		//{"ZAAALLGO", "76561197970672021", "Souls Tyrant"},
		//{"ZAAALLGO", "76561197970672021", "Glimmer of the Sacred Hunt"},
	}

	for _, param := range params {
		status, snaps, err := verifying.Delivery(assetSrc, param.persona, param.steamID, param.item)

		fmt.Println(strings.Repeat("-", 70))
		fmt.Println(fmt.Sprintf("%s -> %s (%s)", param.persona, param.steamID, param.item))
		fmt.Println(strings.Repeat("-", 70))
		fmt.Println("Status:", status)
		if err != nil {
			fmt.Printf("Errored: %s \n\n", err)
			continue
		}

		fmt.Println("Items:", len(snaps))
		if len(snaps) != 0 {
			r := snaps[0]
			fmt.Println("GiftFrom:", r.GiftFrom)
			fmt.Println("DateReceived:", r.DateReceived)
			fmt.Println("Dedication:", r.Dedication)
		}

		fmt.Println("")
	}
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
