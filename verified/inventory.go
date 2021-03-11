package verified

import (
	"fmt"

	"github.com/kudarap/dotagiftx/steam"
)

// InventoryStatus represents inventory status.
type InventoryStatus uint

const (
	InventoryStatusNoHit    InventoryStatus = 10
	InventoryStatusVerified InventoryStatus = 20
	InventoryStatusPrivate  InventoryStatus = 30
	InventoryStatusError    InventoryStatus = 40
)

// Inventory checks item existence on inventory.
//
// Returns an error when request has status error or response body malformed.
func Inventory(source AssetSource, steamID, itemName string) (InventoryStatus, []steam.Asset, error) {
	if steamID == "" || itemName == "" {
		return InventoryStatusError, nil, fmt.Errorf("all params are required")
	}

	assets, err := source(steamID)
	if err != nil {
		if err == steam.ErrInventoryPrivate {
			return InventoryStatusPrivate, nil, nil
		}

		return InventoryStatusError, nil, err
	}

	assets = filterByName(assets, itemName)
	if len(assets) == 0 {
		return InventoryStatusNoHit, assets, nil
	}

	return InventoryStatusVerified, assets, nil
}
