package main

import (
	"fmt"

	"github.com/kudarap/dotagiftx/steaminventory"
)

func main() {
	//status, err := steaminventory.Crawl("76561198088587178")
	//fmt.Println(status, err)

	meta, err := steaminventory.GetMeta("76561198032085715")
	fmt.Println(meta, err)

	inv, err := steaminventory.Get("76561198032085715")
	fmt.Println(inv, err)

	//flat, err := steaminventory.NewFlatInventoryFromV2(*inv)
	//fmt.Println(flat, err)
}
