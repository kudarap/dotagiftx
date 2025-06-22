package main

import (
	"fmt"

	"github.com/kudarap/dotagiftx/phantasm"
)

func main() {
	args := map[string]interface{}{
		"steam_id": "76561198088587178",
	}
	v := phantasm.Main(args)
	fmt.Println(v)
}
