package main

import (
	"fmt"

	"github.com/kudarap/dotagiftx/tools/steaminvcrawler"
)

func main() {
	args := map[string]interface{}{
		"steam_id": "76561198088587178",
	}
	v := steaminvcrawler.Main(args)
	fmt.Println(v)
}
