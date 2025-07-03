package steam

import (
	"context"
	"fmt"
	"time"

	"github.com/kudarap/dotagiftx/filecache"
)

const (
	cacheExpr   = time.Hour * 24
	cachePrefix = "steam"
)

// InventoryAssetWithCache returns a compact format from all
// inventory data with cache.
func InventoryAssetWithCache(ctx context.Context, steamID string) ([]Asset, error) {
	hit, err := filecache.Get(getCacheKey(steamID))
	if err != nil {
		return nil, err
	}
	if hit != nil {
		b, _ := fastjson.Marshal(hit)
		var asset []Asset
		_ = fastjson.Unmarshal(b, &asset)
		return asset, nil
	}

	asset, err := InventoryAsset(ctx, steamID)
	if err != nil {
		return nil, err
	}

	if err = filecache.Set(getCacheKey(steamID), asset, cacheExpr); err != nil {
		return nil, err
	}

	return asset, nil
}

func getCacheKey(steamID string) string {
	return fmt.Sprintf("%s_%s", cachePrefix, steamID)
}
