package main

import (
	"fmt"

	"github.com/kudarap/dotagiftx/steaminvorg"
)

func main() {
	inv, err := steaminvorg.SWR("76561198088587178", true)
	fmt.Println(inv.ToAssets(), err)
}
