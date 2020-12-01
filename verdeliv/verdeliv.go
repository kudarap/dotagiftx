package verdeliv

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

const (
	inventEndpoint  = "https://steamcommunity.com/profiles/%s/inventory/json/570/2?count=5000"
	inventEndpoint2 = "https://steamcommunity.com/inventory/%s/570/2?count=5000"

	steamInventoryAPI = "https://steamcommunity.com/profiles/%s/inventory/json/570/2?start=1090"
	filenameFmt       = "dgx-inventory-%s.json"
)

func Verify(sellerPersona, buyerSteamID, itemName string) ([]flatInventory, error) {
	if sellerPersona == "" || buyerSteamID == "" || itemName == "" {
		return nil, fmt.Errorf("all params are required")
	}

	// Read file cache if exist, else download.
	fp := getSource(buyerSteamID)
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		if err := dlInventory(buyerSteamID); err != nil {
			return nil, fmt.Errorf("could not dl file: %s", err)
		}
	}

	fmt.Println("file", fp)
	res, err := newFlatInventoryFromFile(fp)
	if err != nil {
		return nil, fmt.Errorf("could not parse file: %s", err)
	}

	var fi []flatInventory
	for _, inv := range res {
		if inv.GiftFrom != sellerPersona {
			continue
		}

		if !strings.Contains(strings.Join(inv.Descriptions, "|"), itemName) {
			continue
		}

		fi = append(fi, inv)
	}

	return fi, nil
}

func getSource(steamID string) string {
	return filepath.Join(os.TempDir(), fmt.Sprintf(filenameFmt, steamID))
}

func dlInventory(steamID string) error {
	url := fmt.Sprintf(steamInventoryAPI, steamID)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(getSource(steamID))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
