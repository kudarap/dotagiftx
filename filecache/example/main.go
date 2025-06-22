package main

import (
	"log"
	"time"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/filecache"
)

func main() {
	const testkey = "testkey_123"

	log.Println(filecache.Get(testkey))
	t := time.Now()
	log.Println(filecache.Set(testkey, dotagiftx.Auth{
		ID:        "100id",
		UserID:    "userid299",
		Username:  "akoko",
		CreatedAt: &t,
	}, time.Second*2))
	log.Println(filecache.Get(testkey))
	time.Sleep(time.Second * 3)
	log.Println(filecache.Get(testkey))
}
