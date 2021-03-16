package steaminv

import (
	"fmt"
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

	asset, err := InventoryAsset(steamID)
	if err != nil {
		return nil, err
	}

	if err = cache.Set(getCacheKey(steamID), asset, cacheExpr); err != nil {
		return nil, err
	}

	return asset, nil
}

func getCacheKey(steamID string) string {
	return fmt.Sprintf("%s_%s", cachePrefix, steamID)
}
