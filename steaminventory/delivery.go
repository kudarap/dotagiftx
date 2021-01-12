package steaminventory

import (
	"fmt"
	"strings"
)

func VerifyDelivery(sellerPersona, buyerSteamID, itemName string) ([]flatInventory, error) {
	inv, err := SWR(buyerSteamID)
	if err != nil {
		return nil, fmt.Errorf("could not get inventory: %s", err)
	}

	res, err := NewFlatInventoryFromV2(*inv)
	if err != nil {
		return nil, fmt.Errorf("could not parse file: %s", err)
	}

	var fi []flatInventory
	for _, inv := range res {
		// Checking against seller persona name might not be accurate since
		// buyer can clear gift information that's why it need to snapshot buyer
		// inventory immediately.
		if inv.GiftFrom != sellerPersona {
			continue
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
