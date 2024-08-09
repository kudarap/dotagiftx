package steaminvorg

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	jsoniter "github.com/json-iterator/go"
	localcache "github.com/kudarap/dotagiftx/gokit/cache"
	"github.com/kudarap/dotagiftx/steam"
)

var fastjson = jsoniter.ConfigFastest

const (
	localCacheExpr   = time.Hour
	localCachePrefix = "steaminvorg"
)

// InventoryAssetWithCache returns a compact format from all inventory data with cache.
func InventoryAssetWithCache(steamID string) ([]steam.Asset, error) {
	log.Println("STEAMINVORG CHECK LOCAL CACHE")
	hit, _ := localcache.Get(getCacheKey(steamID))
	if hit != nil {
		log.Println("STEAMINVORG LOCAL CACHE HIT!")

		b, _ := fastjson.Marshal(hit)
		var asset []steam.Asset
		_ = fastjson.Unmarshal(b, &asset)
		return asset, nil
	}

	log.Println("STEAMINVORG NO LOCAL CACHE HIT", steamID)
	asset, err := InventoryAsset(steamID)
	if err != nil {
		log.Println("STEAMINVORG ASSET ERROR", steamID, err)
		return nil, err
	}

	if err = localcache.Set(getCacheKey(steamID), asset, getCacheExpr()); err != nil {
		log.Println("STEAMINVORG LOCAL CACHE SET ERROR", steamID, err)
		return nil, err
	}

	log.Println("STEAMINVORG ASSET DONE", steamID)
	return asset, nil
}

func getCacheKey(steamID string) string {
	return fmt.Sprintf("%s_%s", localCachePrefix, steamID)
}

func getCacheExpr() time.Duration {
	n := 10
	r := rand.Intn(n-(-n)) + (-n)
	d := time.Minute * time.Duration(r)
	return localCacheExpr + d
}
