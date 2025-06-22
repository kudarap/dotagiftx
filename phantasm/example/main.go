package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kudarap/dotagiftx/phantasm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("could not load config: " + err.Error())
	}

	steamID := "76561198088587178"
	if len(os.Args) == 2 {
		steamID = os.Args[1]
	}

	args := map[string]interface{}{
		"steam_id": steamID,
	}
	v := phantasm.Main(args)
	fmt.Println(v)
}
