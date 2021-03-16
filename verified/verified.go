package verified

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
	for _, aa := range a {
		if !strings.Contains(strings.Join(aa.Descriptions, "|"), itemName) &&
			!strings.Contains(aa.Name, itemName) {
			continue
		}
		matches = append(matches, aa)
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
