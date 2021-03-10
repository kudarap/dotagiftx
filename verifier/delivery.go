package verifier

import (
	"fmt"
	"strings"

	"github.com/kudarap/dotagiftx/steam"
)

/*

objective: verify buyer received reserved item from seller

params:
	- seller persona name: check for sender value
	- buyer steam id: for parsing inventory
	- item name: item key name to check against sender

result:
	- detect private inventory
	- detect malformed json
	- support multiple json for large inventory
	- challenge check

process:
	- download json inventory
	- parse json file
	- search item name
	- check sender name

*/

// Delivery checks item existence on buyer's inventory.
//
// Returns an error when request has status error or body malformed.
func Delivery(sellerPersona, buyerSteamID, itemName string) (VerifyStatus, []steam.Asset, error) {
	if sellerPersona == "" || buyerSteamID == "" || itemName == "" {
		return VerifyStatusError, nil, fmt.Errorf("all params are required")
	}

	// Pull inventory data using buyerSteamID.
	assets, err := steam.InventoryAsset(buyerSteamID)
	if err != nil {
		if err == steam.ErrInventoryPrivate {
			return VerifyStatusPrivate, nil, err
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

	for _, ss := range snapshots {
		// Check asset matches the seller persona name.
		//
		// Checking against seller persona name might not be accurate since
		// buyer can clear gift information that's why it need to snapshot
		// buyer inventory immediately.
		if ss.GiftFrom != sellerPersona {
			continue
		}
		status = VerifyStatusSeller
	}

	return status, snapshots, nil
}
