package verifying

import (
	"fmt"

	dgx "github.com/kudarap/dotagiftx"
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
func Delivery(
	source AssetSource,
	sellerPersona,
	buyerSteamID,
	itemName string,
) (dgx.DeliveryStatus, []dgx.SteamAsset, error) {
	if sellerPersona == "" || buyerSteamID == "" || itemName == "" {
		return dgx.DeliveryStatusError, nil, fmt.Errorf("all params are required")
	}

	// Pull inventory data using buyerSteamID.
	assets, err := source(buyerSteamID)
	if err != nil {
		if err == steam.ErrInventoryPrivate {
			return dgx.DeliveryStatusPrivate, nil, nil
		}

		return dgx.DeliveryStatusError, nil, err
	}

	assets = filterByName(assets, itemName)
	if len(assets) == 0 {
		return dgx.DeliveryStatusNoHit, assets, nil
	}

	status := dgx.DeliveryStatusNameVerified
	// Check asset sender matches the seller persona name.
	//
	// NOTE! checking against seller persona name might not be accurate since
	// buyer can clear gift information that's why it need to snapshot
	// buyer inventory immediately.
	for _, ss := range assets {
		if ss.GiftFrom != sellerPersona {
			continue
		}
		status = dgx.DeliveryStatusSenderVerified
	}

	return status, assets, nil
}
