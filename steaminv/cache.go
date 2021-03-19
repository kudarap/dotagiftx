package steaminv

import (
	"fmt"
	"math/rand"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/kudarap/dotagiftx/gokit/cache"
	"github.com/kudarap/dotagiftx/steam"
)

var fastjson = jsoniter.ConfigFastest

const (
	cacheExpr   = time.Hour * 24
	cachePrefix = "steaminv"
)

// InventoryAsset returns a compact format from all
// inventory data with cache.
func InventoryAssetWithCache(steamID string) ([]steam.Asset, error) {
	hit, _ := cache.Get(getCacheKey(steamID))
	if hit != nil {
		b, _ := fastjson.Marshal(hit)
		var asset []steam.Asset
		_ = fastjson.Unmarshal(b, &asset)
		return asset, nil
	}

	fmt.Println("STEAMINV NO LOCAL CACHE HIT", steamID)
	asset, err := InventoryAsset(steamID)
	if err != nil {
		return nil, err
	}

	if err = cache.Set(getCacheKey(steamID), asset, getCacheExpr()); err != nil {
		return nil, err
	}

	return asset, nil
}

func getCacheKey(steamID string) string {
	return fmt.Sprintf("%s_%s", cachePrefix, steamID)
}

func init() {
	rand.Seed(time.Now().Unix())
}

func getCacheExpr() time.Duration {
	n := 10
	r := rand.Intn(n-(-n)) + (-n)
	fmt.Println(r)
	d := time.Minute * time.Duration(r)
	fmt.Println(cacheExpr + d)
	return cacheExpr + d
}
