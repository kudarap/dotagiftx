package verified

import (
	"strings"

	"github.com/kudarap/dotagiftx/steam"
)

// AssetSource represents inventory asset source provider.
type AssetSource func(steamID string) ([]steam.Asset, error)

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
