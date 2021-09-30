package steaminv

import (
	"fmt"
	"math/rand"
	"time"

	jsoniter "github.com/json-iterator/go"
	localcache "github.com/kudarap/dotagiftx/gokit/cache"
	"github.com/kudarap/dotagiftx/steam"
)

var fastjson = jsoniter.ConfigFastest

const (
	// localcacheExpr   = time.Hour * 24
	localcacheExpr   = time.Minute * 10
	localcachePrefix = "steaminv"
)

// InventoryAsset returns a compact format from all
// inventory data with cache.
func InventoryAssetWithCache(steamID string) ([]steam.Asset, error) {
	hit, _ := localcache.Get(getCacheKey(steamID))
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

	if err = localcache.Set(getCacheKey(steamID), asset, getCacheExpr()); err != nil {
		return nil, err
	}

	return asset, nil
}

func getCacheKey(steamID string) string {
	return fmt.Sprintf("%s_%s", localcachePrefix, steamID)
}

func init() {
	rand.Seed(time.Now().Unix())
}

func getCacheExpr() time.Duration {
	n := 10
	r := rand.Intn(n-(-n)) + (-n)
	d := time.Minute * time.Duration(r)
	return localcacheExpr + d
}
