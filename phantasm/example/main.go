package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/joho/godotenv"
	"github.com/kudarap/dotagiftx/phantasm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("could not load config: " + err.Error())
	}

	var steamID string
	var precheck bool
	flag.StringVar(&steamID, "steam_id", "76561198088587178", "")
	flag.BoolVar(&precheck, "precheck", false, "")
	flag.Parse()

	args := map[string]interface{}{"steam_id": steamID}
	if precheck {
		args["precheck"] = true
	}
	v := phantasm.Main(args)
	b, _ := json.MarshalIndent(v, "", "\t")
	os.Stdout.Write(b)
}
