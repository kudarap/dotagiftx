package steaminvorg

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/kudarap/dotagiftx/file/filecache"
	"github.com/kudarap/dotagiftx/steam"
)

var fastjson = jsoniter.ConfigFastest

const (
	localCacheExpr   = time.Hour
	localCachePrefix = "steaminvorg"
)

// InventoryAssetWithCache returns a compact format from all inventory data with cache.
func InventoryAssetWithCache(ctx context.Context, steamID string) ([]steam.Asset, error) {
	sharedLogger.Info("check local cache", "steam_id", steamID)
	hit, _ := filecache.Get(getCacheKey(steamID))
	if hit != nil {
		sharedLogger.Info("cache hit", "steam_id", steamID)

		b, _ := fastjson.Marshal(hit)
		var asset []steam.Asset
		_ = fastjson.Unmarshal(b, &asset)
		return asset, nil
	}

	sharedLogger.Info("no local cache hit", "steam_id", steamID)
	asset, err := InventoryAsset(ctx, steamID)
	if err != nil {
		sharedLogger.Error("asset error", "steam_id", steamID, "err", err)
		return nil, err
	}

	if err = filecache.Set(getCacheKey(steamID), asset, getCacheExpr()); err != nil {
		sharedLogger.Error("local cache set error", "steam_id", steamID, "err", err)
		return nil, err
	}

	sharedLogger.Info("asset done", "steam_id", steamID)
	return asset, nil
}

func getCacheKey(steamID string) string {
	return fmt.Sprintf("%s_%s", localCachePrefix, steamID)
}

func getCacheExpr() time.Duration {
	const jitter = 10
	n, _ := rand.Int(rand.Reader, big.NewInt(2*jitter))
	r := int(n.Int64()) - jitter
	d := time.Minute * time.Duration(r)
	return localCacheExpr + d
}
