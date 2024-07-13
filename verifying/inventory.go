package verifying

import (
	"fmt"

	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/steam"
)

// Inventory checks item existence on inventory.
//
// Returns an error when request has status error or response body malformed.
func Inventory(source AssetSource, steamID, itemName string) (dgx.InventoryStatus, []steam.Asset, error) {
	if steamID == "" || itemName == "" {
		return dgx.InventoryStatusError, nil, fmt.Errorf("all params are required")
	}

	assets, err := source(steamID)
	if err != nil {
		if err == steam.ErrInventoryPrivate {
			return dgx.InventoryStatusPrivate, nil, nil
		}

		return dgx.InventoryStatusError, nil, err
	}

	assets = filterByName(assets, itemName)
	assets = filterByGiftable(assets)
	if len(assets) == 0 {
		return dgx.InventoryStatusNoHit, assets, nil
	}

	return dgx.InventoryStatusVerified, assets, nil
}
