package verified

import (
	"fmt"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/steam"
)

// Inventory checks item existence on inventory.
//
// Returns an error when request has status error or response body malformed.
func Inventory(source AssetSource, steamID, itemName string) (core.InventoryStatus, []steam.Asset, error) {
	if steamID == "" || itemName == "" {
		return core.InventoryStatusError, nil, fmt.Errorf("all params are required")
	}

	assets, err := source(steamID)
	if err != nil {
		if err == steam.ErrInventoryPrivate {
			return core.InventoryStatusPrivate, nil, nil
		}

		return core.InventoryStatusError, nil, err
	}

	assets = filterByName(assets, itemName)
	assets = filterByGiftable(assets)
	if len(assets) == 0 {
		return core.InventoryStatusNoHit, assets, nil
	}

	return core.InventoryStatusVerified, assets, nil
}
