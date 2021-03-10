package steaminv

import (
	"fmt"
	"strings"

	"github.com/kudarap/dotagiftx/steam"
)

func VerifyDelivery(sellerPersona, buyerSteamID, itemName string) ([]steam.Asset, error) {
	inv, err := SWR(buyerSteamID)
	if err != nil {
		return nil, fmt.Errorf("could not get inventory: %s", err)
	}
	if inv == nil {
		return nil, fmt.Errorf("inventory empty result")
	}

	var fi []steam.Asset
	for _, inv := range inv.ToAssets() {
		// Checking against seller persona name might not be accurate since
		// buyer can clear gift information that's why it need to snapshot buyer
		// inventory immediately.
		if inv.GiftFrom != sellerPersona {
			//continue
		}

		// Checks target item name from description and name field.
		if !strings.Contains(strings.Join(inv.Descriptions, "|"), itemName) &&
			!strings.Contains(inv.Name, itemName) {
			continue
		}

		fi = append(fi, inv)
	}

	return fi, nil
}
