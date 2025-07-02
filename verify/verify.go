package verifying

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/kudarap/dotagiftx/steam"
)

// AssetSource represents inventory asset source provider.
type AssetSource func(steamID string) ([]steam.Asset, error)

type AssetSourceProvider func(steamID string) (provider string, sa []steam.Asset, err error)

type Source struct {
}

func MultiAssetSource(providers map[string]AssetSource) AssetSource {
	logger := slog.Default()
	return func(steamID string) ([]steam.Asset, error) {
		for name, source := range providers {
			logger.Info("attempting to fetch asset",
				"provider", name,
				"steam_id", steamID,
			)
			assets, err := source(steamID)
			if err != nil {
				logger.Error("failed to fetch assets",
					"provider", name, "steam_id", steamID,
					"err", err,
				)
				continue
			}
			return assets, nil
		}
		return nil, fmt.Errorf("all asset providers attempted: %s", steamID)
	}
}

func MultiAssetSourceProvider(providers map[string]AssetSource) AssetSourceProvider {
	logger := slog.Default()
	return func(steamID string) (string, []steam.Asset, error) {
		for name, source := range providers {
			assets, err := source(steamID)
			if err != nil {
				logger.Error("failed to fetch assets",
					"provider", name, "steam_id", steamID,
					"err", err,
				)
				continue
			}
			return name, assets, nil
		}

		return "", nil, fmt.Errorf("all asset providers attempted: %s", steamID)
	}
}

// filterByName filters item that matches the name or in the description
// that supports unbundled items.
func filterByName(a []steam.Asset, itemName string) []steam.Asset {
	var matches []steam.Asset
	for _, asset := range a {
		// Strip "bundle" string to cover items that unbundled:
		// - Dipper the Destroyer Bundle
		// - The Abscesserator Bundle
		itemName = strings.TrimSpace(strings.TrimSuffix(itemName, "Bundle"))
		if !strings.Contains(strings.Join(asset.Descriptions, "|"), itemName) &&
			!strings.Contains(asset.Name, itemName) {
			continue
		}

		// Excluded golden variant of the item.
		if asset.IsGoldenVariant(itemName) {
			continue
		}

		matches = append(matches, asset)
	}
	return matches
}

// filterByGiftable filters item that can be gifted and allowed to be listed.
func filterByGiftable(a []steam.Asset) []steam.Asset {
	var matches []steam.Asset
	for _, aa := range a {
		// Is the item unbundled but giftable?
		//
		// Is the item immortal and does not say its giftable?
		// This fixes the removed "Gift once" string on description
		// recently by Valve.
		if !aa.GiftOnce && (aa.IsCollectorsCache() || !aa.IsImmortal()) {
			continue
		}

		matches = append(matches, aa)
	}
	return matches
}
