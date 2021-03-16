package main

import (
	"log"
	"time"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/steam/cache"
)

func main() {
	log.Println(cache.Get("123412415"))
	t := time.Now()
	log.Println(cache.Set("123412415", core.Auth{
		ID:        "100id",
		UserID:    "userid299",
		Username:  "akoko",
		CreatedAt: &t,
	}, time.Second*2))
	log.Println(cache.Get("123412415"))
	time.Sleep(time.Second * 3)
	log.Println(cache.Get("123412415"))

}
