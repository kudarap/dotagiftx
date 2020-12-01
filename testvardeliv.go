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
		{"kudarap", "76561198042690669", "Riddle of the Hierophant"},
		{"kudarap", "76561198436826874", "Fowl Omen"},
	}

	for _, param := range params {
		res, err := verdeliv.Verify(param.persona, param.steamID, param.item)
		if err != nil {
			fmt.Println("could not verify!", err)
			return
		}

		fmt.Println(fmt.Sprintf("%s check %s hit: %d", param.persona, param.item, len(res)))
		for _, rr := range res {
			fmt.Println(rr.Name)
		}

		fmt.Println("")
	}

}
