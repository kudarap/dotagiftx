package steam

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var fastjson = jsoniter.ConfigFastest

const Dota2AppID = 570
const inventoryEndpoint = "https://steamcommunity.com/profiles/%s/inventory/json/%d/2"

func reqDota2Inventory(steamID string) (*http.Response, error) {
	url := fmt.Sprintf(inventoryEndpoint, steamID, Dota2AppID)
	return http.Get(url)
}

// Asset represents compact inventory base of RawInventory model.
type Asset struct {
	AssetID      string   `json:"asset_id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Type         string   `json:"type"`
	Hero         string   `json:"hero"`
	GiftFrom     string   `json:"gift_from"`
	DateReceived string   `json:"date_received"`
	Dedication   string   `json:"dedication"`
	GiftOnce     bool     `json:"gift_once"`
	NotTradable  bool     `json:"not_tradable"`
	Descriptions []string `json:"descriptions"`
}

func InventoryAsset(steamID string) ([]Asset, error) {
	r, err := reqDota2Inventory(steamID)
	if err != nil {
		return nil, fmt.Errorf("could send request: %s", err)
	}
	defer r.Body.Close()
	return assetParser(r.Body)
}

var ErrInventoryPrivate = errors.New("profile inventory is private")

func assetParser(r io.Reader) ([]Asset, error) {
	raw, err := inventoryParser(r)
	if err != nil {
		return nil, err
	}
	if raw.IsPrivate() {
		return nil, ErrInventoryPrivate
	}
	if raw.Error != "" {
		return nil, fmt.Errorf(raw.Error)
	}

	return raw.ToAssets(), nil
}

// RawInventory represents steam's raw inventory data model.
type RawInventory struct {
	Success      bool                         `json:"success"`
	More         bool                         `json:"more"`
	MoreStart    RawInventoryPageOffset       `json:"more_start"`
	Assets       map[string]RawInventoryAsset `json:"rgInventory"`
	Descriptions map[string]RawInventoryDesc  `json:"rgDescriptions"`
	Error        string                       `json:"Error"`
}

func (i RawInventory) IsPrivate() bool {
	return strings.ToUpper(i.Error) == "THIS PROFILE IS PRIVATE."
}

func (i *RawInventory) ToAssets() []Asset {
	// Collate asset map ids for fast inventory asset id look up.
	assetMapIDs := map[string]string{}
	for _, aa := range i.Assets {
		assetMapIDs[fmt.Sprintf("%s_%s", aa.ClassID, aa.InstanceID)] = aa.ID
	}

	// Composes and collect inventory on flat format.
	var assets []Asset
	for ci, ii := range i.Descriptions {
		a := ii.toAsset()
		a.AssetID = assetMapIDs[ci]
		assets = append(assets, a)
	}

	return assets
}

// Inventory retrieve data from API and parse into RawInventory.
func Inventory(steamID string) (*RawInventory, error) {
	r, err := reqDota2Inventory(steamID)
	if err != nil {
		return nil, fmt.Errorf("could send request: %s", err)
	}
	defer r.Body.Close()
	return inventoryParser(r.Body)
}

func inventoryParser(r io.Reader) (*RawInventory, error) {
	raw := &RawInventory{}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if err = fastjson.Unmarshal(b, raw); err != nil {
		return nil, err
	}
	return raw, nil
}

// RawInventoryAsset represents steam's raw asset inventory data model.
type RawInventoryAsset struct {
	ID         string `json:"id"`
	ClassID    string `json:"classid"`
	InstanceID string `json:"instanceid"`
}

// Inventory description field prefix and flags.
const (
	inventPrefixHero         = "Used By: "
	inventPrefixGiftFrom     = "Gift From: "
	inventPrefixDateReceived = "Date Received: "
	inventPrefixDedication   = "Dedication: "
	inventFlagGiftOnce       = "( Not Tradable )"
	inventFlagNotTradable    = "( This item may be gifted once )"
)

// RawInventoryDesc represents steam's raw description inventory data model.
type RawInventoryDesc struct {
	ClassID      string                  `json:"classid"`
	InstanceID   string                  `json:"instanceid"`
	Name         string                  `json:"name"`
	Image        string                  `json:"icon_url_large"`
	Type         string                  `json:"type"`
	Descriptions RawInventoryItemDetails `json:"descriptions"`
}

func (d RawInventoryDesc) toAsset() Asset {
	fi := Asset{
		Name:  d.Name,
		Image: d.Image,
		Type:  d.Type,
	}

	var desc []string
	for _, dd := range d.Descriptions {
		v := dd.Value
		desc = append(desc, v)
		if pv, ok := extractValueFromPrefix(v, inventPrefixHero); ok {
			fi.Hero = pv
		}
		if pv, ok := extractValueFromPrefix(v, inventPrefixGiftFrom); ok {
			fi.GiftFrom = pv
		}
		if pv, ok := extractValueFromPrefix(v, inventPrefixDateReceived); ok {
			fi.DateReceived = pv
		}
		if pv, ok := extractValueFromPrefix(v, inventPrefixDedication); ok {
			fi.Dedication = pv
		}
		if isFlagExists(v, inventFlagGiftOnce) {
			fi.GiftOnce = true
		}
		if isFlagExists(v, inventFlagNotTradable) {
			fi.NotTradable = true
		}
	}
	fi.Descriptions = desc

	return fi
}

// RawInventoryItemDetails represents steam's raw description detail values data model.
type RawInventoryItemDetails []struct {
	Value string `json:"value"`
}

func (d *RawInventoryItemDetails) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == `""` {
		*d = nil
		return nil
	}

	var details []struct {
		Value string `json:"value"`
	}
	if err := fastjson.Unmarshal(data, &details); err != nil {
		return err
	}
	*d = details
	return nil
}

type RawInventoryPageOffset int

func (po *RawInventoryPageOffset) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == `false` {
		*po = 0
		return nil
	}

	o := 0
	if err := fastjson.Unmarshal(data, &o); err != nil {
		return err
	}
	*po = RawInventoryPageOffset(o)
	return nil
}

func extractValueFromPrefix(s, prefix string) (value string, ok bool) {
	if !strings.HasPrefix(strings.ToUpper(s), strings.ToUpper(prefix)) {
		return
	}

	return strings.TrimPrefix(s, prefix), true
}

func isFlagExists(s, flag string) (ok bool) {
	return strings.ToUpper(s) == strings.ToUpper(flag)
}
