package verdeliv

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/headzoo/surf"
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
	inventEndpoint  = "https://steamcommunity.com/profiles/%s/inventory/json/570/2"
	inventEndpoint2 = "https://steamcommunity.com/inventory/%s/570/2"

	//steamInventoryAPI = "https://steamcommunity.com/profiles/%s/inventory/json/570/2?start=963"
	steamInventoryAPI = "https://steamcommunity.com/profiles/%s/inventory/json/570/2"
	filenameFmt       = "dgx-inv-%s.json"
)

func Verify(sellerPersona, buyerSteamID, itemName string) ([]flatInventory, error) {
	if sellerPersona == "" || buyerSteamID == "" || itemName == "" {
		return nil, fmt.Errorf("all params are required")
	}

	// Read file cache if exist, else download.
	fp := getSource(buyerSteamID)
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		// Sleep for the next request
		time.Sleep(time.Minute * 5)
		if err := dlInventory(buyerSteamID); err != nil {
			return nil, fmt.Errorf("could not dl file: %s", err)
		}
	}

	res, err := newFlatInventoryFromFile(fp)
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

func getSource(steamID string) string {
	return filepath.Join(os.TempDir(), fmt.Sprintf(filenameFmt, steamID))
}

func dlInventory(steamID string) error {
	url := fmt.Sprintf(steamInventoryAPI, steamID)
	fmt.Println("downloading", url, "...")

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	inv := struct {
		Success bool `json:"success"`
	}{}
	data, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(data, &inv); err != nil {
		return err
	}
	// Skip caching if no success.
	if !inv.Success || string(data) == "" {
		return fmt.Errorf("please try again later: %s", data)
	}

	out, err := os.Create(getSource(steamID))
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = out.Write(data)
	return err
}

func dlInventory2(steamID string) error {
	url := fmt.Sprintf(steamInventoryAPI, steamID)
	fmt.Println("downloading", url, "...")

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

	inv := struct {
		Success bool `json:"success"`
	}{}
	data, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(data, &inv); err != nil {
		return err
	}
	// Skip caching if no success.
	if !inv.Success {
		return fmt.Errorf("please try again later")
	}

	_, err = io.Copy(out, resp.Body)
	return err
}

func dlInventorySurf(steamID string) error {
	url := fmt.Sprintf(steamInventoryAPI, steamID)
	fmt.Println("downloading", url, "with surf...")

	bow := surf.NewBrowser()
	if err := bow.Open(url); err != nil {
		panic(err)
	}
	respStr := html.UnescapeString(bow.Body())

	f, err := os.Create(getSource(steamID))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(respStr)
	return nil
}
