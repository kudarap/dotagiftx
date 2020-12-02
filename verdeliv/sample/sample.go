package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/verdeliv"
)

func main() {
	fmt.Println("TESTING VARDELIV")

	params := []struct {
		persona, steamID, item string
	}{
		{"karosu!", "76561198088587178", "Ravenous Abyss"},
		{"Araragi-", "76561198809365008", "Master of the Searing Path"},
		{"Accel", "76561198042690669", "Forsworn Legacy"},
		{"Dark Knight", "76561198116319576", "Legends of Darkheart Pursuit"},

		//{"ZAAALLGO", "76561197970672021", "Souls Tyrant"},
		//{"ZAAALLGO", "76561197970672021", "Glimmer of the Sacred Hunt"},
		{"ZAAALLGO", "76561198203634725", "Cunning Corsair"},

		{"Berserk", "76561198355627060", "Shattered Greatsword"},
		{"Berserk", "76561198042690669", "Ancient Inheritance"},
		{"Berserk", "76561198042690669", "Poacher's Bane"},
		{"Berserk", "76561198042690669", "Allure of the Faeshade Flower"},
		{"Berserk", "76561198256569879", "Endless Night"},
		{"Berserk", "76561198139657329", "Glimmer of the Sacred Hunt"},

		{"kudarap", "76561198042690669", "Riddle of the Hierophant"},
		{"kudarap", "76561198436826874", "Fowl Omen"},
		{"kudarap", "76561198170463425", "Cinder Sensei"},
		{"kudarap", "76561198872556187", "Grim Destiny"},
		{"kudarap", "76561198042690669", "Scorched Amber"},
		{"kudarap", "76561198872556187", "Tales of the Windward Rogue"},
		{"kudarap", "76561198203634725", "Adornments of the Jade Emissary"},
		{"gippeum", "76561198088587178", "Elements of the Endless Plane"},
	}

	for _, param := range params {
		res, err := verdeliv.Verify(param.persona, param.steamID, param.item)
		if err != nil {
			fmt.Println("could not verify!", err)
			continue
		}

		fmt.Println(strings.Repeat("-", 70))
		fmt.Println(fmt.Sprintf("%s -> %s (%s)", param.persona, param.steamID, param.item))
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
