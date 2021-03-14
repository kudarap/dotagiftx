package main

import (
	"log"
	"time"

	"github.com/kudarap/dotagiftx/core"

	"github.com/kudarap/dotagiftx/steam/cache"
)

func main() {
	t := time.Now()
	log.Println(cache.Set("test100", core.Auth{
		ID:        "100id",
		UserID:    "userid299",
		Username:  "akoko",
		CreatedAt: &t,
	}, time.Second*2))
	log.Println(cache.Get("test100"))
	time.Sleep(time.Second * 3)
	log.Println(cache.Get("test100"))

}
