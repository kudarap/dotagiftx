package verifying

import (
	"strings"

	"github.com/kudarap/dotagiftx/steam"
)

// AssetSource represents inventory asset source provider.
type AssetSource func(steamID string) ([]steam.Asset, error)

// filterByName filters item that matches the name or in the description
// that supports un-bundled items.
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
		// Is the item is unbundled but giftable?
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
