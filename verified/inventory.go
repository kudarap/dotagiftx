package verified

import (
	"fmt"
	"strings"

	"github.com/kudarap/dotagiftx/steam"
)

func Inventory(source AssetSource, steamID, itemName string) (VerifyStatus, []steam.Asset, error) {
	if steamID == "" || itemName == "" {
		return VerifyStatusError, nil, fmt.Errorf("all params are required")
	}

	// Pull inventory data using buyerSteamID.
	assets, err := source(steamID)
	if err != nil {
		if err == steam.ErrInventoryPrivate {
			return VerifyStatusPrivate, nil, nil
		}

		return VerifyStatusError, nil, err
	}

	status := VerifyStatusNoHit

	// Check asset existence base on item name.
	var snapshots []steam.Asset
	for _, asset := range assets {
		if !strings.Contains(strings.Join(asset.Descriptions, "|"), itemName) &&
			!strings.Contains(asset.Name, itemName) {
			continue
		}
		snapshots = append(snapshots, asset)
		status = VerifyStatusItem
	}

	return status, snapshots, nil
}
