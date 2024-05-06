package verified

import (
	"fmt"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/steam"
)

// Inventory checks item existence on inventory.
//
// Returns an error when request has status error or response body malformed.
func Inventory(source AssetSource, steamID, itemName string) (dotagiftx.InventoryStatus, []steam.Asset, error) {
	if steamID == "" || itemName == "" {
		return dotagiftx.InventoryStatusError, nil, fmt.Errorf("all params are required")
	}

	assets, err := source(steamID)
	if err != nil {
		if err == steam.ErrInventoryPrivate {
			return dotagiftx.InventoryStatusPrivate, nil, nil
		}

		return dotagiftx.InventoryStatusError, nil, err
	}

	assets = filterByName(assets, itemName)
	assets = filterByGiftable(assets)
	if len(assets) == 0 {
		return dotagiftx.InventoryStatusNoHit, assets, nil
	}

	return dotagiftx.InventoryStatusVerified, assets, nil
}
