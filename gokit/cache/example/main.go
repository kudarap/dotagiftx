package main

import (
	"log"
	"time"

	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/gokit/cache"
)

func main() {
	const testkey = "testkey_123"

	log.Println(cache.Get(testkey))
	t := time.Now()
	log.Println(cache.Set(testkey, dgx.Auth{
		ID:        "100id",
		UserID:    "userid299",
		Username:  "akoko",
		CreatedAt: &t,
	}, time.Second*2))
	log.Println(cache.Get(testkey))
	time.Sleep(time.Second * 3)
	log.Println(cache.Get(testkey))

}
