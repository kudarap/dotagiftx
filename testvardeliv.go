package main

import (
	"fmt"

	"github.com/kudarap/dotagiftx/verdeliv"
)

func main() {
	fmt.Println("TESTING VARDELIV")

	params := []struct {
		persona, steamID, item string
	}{
		{"karosu!", "76561198088587178", "Ravenous Abyss"},
		{"ZAAALLGO", "76561198203634725", "Souls Tyrant"},
		{"Araragi-", "76561198809365008", "Master of the Searing Path"},

		{"Berserk", "76561198355627060", "Shattered Greatsword"},
		{"Berserk", "76561198042690669", "Ancient Inheritance"},
		{"Berserk", "76561198042690669", "Poacher's Bane"},

		{"kudarap", "76561198042690669", "Riddle of the Hierophant"},
		{"kudarap", "76561198436826874", "Fowl Omen"},
		{"kudarap", "76561198170463425", "Cinder Sensei"},
		{"kudarap", "76561198872556187", "Grim Destiny"},
		// Chiw
		{"Berserk", "76561198042690669", "Allure of the Faeshade Flower"},
		{"Berserk", "76561198256569879", "Endless Night"},
		{"Berserk", "76561198139657329", "Glimmer of the Sacred Hunt"},
		{"Dark Knight", "76561198116319576", "Legends of Darkheart Pursuit"},
		{"kudarap", "76561198042690669", "Scorched Amber"},
		{"kudarap", "76561198872556187", "Tales of the Windward Rogue"},
	}

	for _, param := range params {
		res, err := verdeliv.Verify(param.persona, param.steamID, param.item)
		if err != nil {
			fmt.Println("could not verify!", err)
			continue
		}

		fmt.Println(fmt.Sprintf("%s check %s hit: %d", param.persona, param.item, len(res)))
		for _, rr := range res {
			fmt.Println(rr.Name)
		}

		fmt.Println("")
	}

}
