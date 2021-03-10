package main

import (
	"fmt"

	"github.com/kudarap/dotagiftx/steaminv"
)

func main() {
	inv, err := steaminv.SWR("76561198088587178")
	fmt.Println(inv.ToAssets(), err)
}
