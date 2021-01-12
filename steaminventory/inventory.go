package steaminventory

import (
	"fmt"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var fastjson = jsoniter.ConfigFastest

// Inventory description field prefix and flags.
const (
	inventPrefixHero         = "Used By: "
	inventPrefixGiftFrom     = "Gift From: "
	inventPrefixDateReceived = "Date Received: "
	inventPrefixDedication   = "Dedication: "
	inventFlagGiftOnce       = "( Not Tradable )"
	inventFlagNotTradable    = "( This item may be gifted once )"
)

type (
	// allInventory represents aggregated steam inventory.
	allInventory struct {
		Assets       []asset                `json:"allInventory"`
		Descriptions map[string]description `json:"allDescriptions"`
	}

	// inventory represents steam's raw inventory data model.
	inventory struct {
		Success      bool                   `json:"success"`
		More         bool                   `json:"more"`
		MoreStart    paginationOffset       `json:"more_start"`
		Assets       map[string]asset       `json:"rgInventory"`
		Descriptions map[string]description `json:"rgDescriptions"`
		Error        string                 `json:"Error"`
	}

	// asset represents steam's raw asset inventory data model.
	asset struct {
		ID         string `json:"assetid"`
		ClassID    string `json:"classid"`
		InstanceID string `json:"instanceid"`
	}

	// description represents steam's raw description inventory data model.
	description struct {
		ClassID      string      `json:"classid"`
		InstanceID   string      `json:"instanceid"`
		Name         string      `json:"name"`
		Image        string      `json:"icon_url_large"`
		Type         string      `json:"type"`
		Descriptions itemDetails `json:"descriptions"`
	}

	// itemDetails represents steam's raw description detail values data model.
	itemDetails []struct {
		Value string `json:"value"`
	}

	paginationOffset int

	// flatInventory represents a flat formatted inventory base of steam model.
	flatInventory struct {
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

	inventorySource interface {
		Get(steamID string)
	}
)

func ToFlatFormat(inv allInventory) ([]flatInventory, error) {
	// Collate asset map ids for fast inventory asset id look up.
	assetIDs := map[string]string{}
	for _, aa := range inv.Assets {
		assetIDs[fmt.Sprintf("%s_%s", aa.ClassID, aa.InstanceID)] = aa.ID
	}

	// Composes and collect inventory on flat format.
	var flat []flatInventory
	for ci, ii := range inv.Descriptions {
		fi := ii.toFlatInventory()
		fi.AssetID = assetIDs[ci]
		flat = append(flat, fi)
	}

	return flat, nil
}

func (d description) toFlatInventory() flatInventory {
	fi := flatInventory{
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

func (d *itemDetails) UnmarshalJSON(data []byte) error {
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
	*d = itemDetails(details)
	return nil
}

func (po *paginationOffset) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == `false` {
		*po = 0
		return nil
	}

	o := 0
	if err := fastjson.Unmarshal(data, &o); err != nil {
		return err
	}
	*po = paginationOffset(o)
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
