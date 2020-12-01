package verdeliv

import (
	"fmt"
	"io/ioutil"

	jsoniter "github.com/json-iterator/go"
)

var fastjson = jsoniter.ConfigFastest

/*

Take inventory json files input
	- support multiple json for paginated large inventory

Parse steam inventory
	- detect malformed json
	- detect private inventory
	- total number of items
	- pagination

NOTES!
	- use asset_id to provide url for specific item with this format
		- https://steamcommunity.com/profiles/{steam_id}/inventory/#570_2_{asset_id}
	- to get asset_id use description classid and instanceid and look from the assets map
*/

const (
	inventEndpoint  = "https://steamcommunity.com/profiles/%s/inventory/json/570/2?count=5000"
	inventEndpoint2 = "https://steamcommunity.com/inventory/%s/570/2?count=5000"
)

type (
	inventory struct {
		Success      bool                   `json:"success"`
		More         bool                   `json:"more"`
		MoreStart    bool                   `json:"more_start"`
		Assets       map[string]asset       `json:"rgInventory"`
		Descriptions map[string]description `json:"rgDescriptions"`
	}

	asset struct {
		ID         string `json:"id"`
		ClassID    string `json:"classid"`
		InstanceID string `json:"instanceid"`
	}

	description struct {
		ClassID      string        `json:"classid"`
		InstanceID   string        `json:"instanceid"`
		Name         string        `json:"name"`
		Image        string        `json:"icon_url_large"`
		Type         string        `json:"type"`
		Descriptions []itemDetails `json:"descriptions"`
	}

	itemDetails struct {
		Value string `json:"value"`
	}
)

func newInventoryFromFile(path string) (*inventory, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %s", err)
	}

	inv := &inventory{}
	if err := fastjson.Unmarshal(data, inv); err != nil {
		return nil, fmt.Errorf("could not parse json: %s", err)
	}

	fmt.Println("parsed", len(inv.Descriptions))
	return inv, nil
}
