package verifier

import (
	"github.com/kudarap/dotagiftx/steam"
	"github.com/kudarap/dotagiftx/steam/inventory"
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

// VerifyDelivery checks item existence on buyer's inventory.
//
// Returns an error when request has status error or body malformed.
func VerifyDelivery(sellerPersona, buyerSteamID, itemName string) (VerifyStatus, []Asset, error) {
	// validate params

	// pull inventory data from buyerSteamID
	inventory.FromFile()

	steam.Inventory(buyerSteamID)

	// check asset existence base on item name

	// check asset matches the seller persona name

	return 0, nil, nil
}
