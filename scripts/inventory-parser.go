package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var fastjson = jsoniter.ConfigFastest

type Item struct {
	ClassID     string      `json:"classid"`
	Name        string      `json:"name"`
	Image       string      `json:"icon_url_large"`
	Type        string      `json:"type"`
	DescRaw     interface{} `json:"descriptions"`
	Description string
	Hero        string
}

func (i Item) stringifyDesc() (description, hero string) {
	s := fmt.Sprintf("%s", i.DescRaw)

	// Extract hero name
	const seg1 = "value:Used By: "
	for _, hh := range strings.Split(s, "] map[") {
		if strings.Contains(hh, seg1) {
			hs := strings.Split(hh, seg1)
			hero = hs[len(hs)-1]
			break
		}

	}

	return s, hero
}

type Inventory struct {
	Items map[string]Item `json:"rgDescriptions"`
}

const (
	steamCDN          = "https://steamcommunity-a.akamaihd.net/economy/image/"
	steamInventoryAPI = "https://steamcommunity.com/profiles/%s/inventory/json/570/2?"

	filenameFmt = "dgx-inventory-%s.json"
)

func main() {
	filterDescPtr := flag.String("filter", "2019", "description filter")
	filterTypePtr := flag.String("type", "Mythical Bundle", "type filter")
	filterNamePtr := flag.String("name", "", "name filter")
	flag.Parse()

	ids := flag.Args()
	if len(ids) == 0 {
		fmt.Println("steam id required")
		return
	}
	steamID := ids[0]

	// Read file cache if exist, else download
	cacheSrc := getSource(steamID)
	if _, err := os.Stat(cacheSrc); os.IsNotExist(err) {
		if err := dlCache(steamID); err != nil {
			fmt.Println("could not dl cache", err)
			return
		}
	}

	data, err := ioutil.ReadFile(cacheSrc)
	if err != nil {
		fmt.Println("could not read file:", err)
		return
	}

	var inv Inventory
	if err := fastjson.Unmarshal(data, &inv); err != nil {
		fmt.Println("could not parse json:", err)
	}
	fmt.Println("parsed", len(inv.Items))

	const filter = "International 2019"
	items := map[string]Item{}
	for _, ii := range inv.Items {
		if !strings.Contains(strings.ToLower(ii.Name), strings.ToLower(*filterNamePtr)) {
			continue
		}
		if !strings.Contains(strings.ToLower(ii.Type), strings.ToLower(*filterTypePtr)) {
			continue
		}

		desc, hero := ii.stringifyDesc()
		if !strings.Contains(strings.ToLower(desc), strings.ToLower(*filterDescPtr)) ||
			hero == "" {
			continue
		}

		items[ii.Name] = Item{
			ClassID:     ii.ClassID,
			Name:        ii.Name,
			Type:        ii.Type,
			Image:       ii.Image,
			Description: desc,
			Hero:        hero,
		}
	}

	for _, ii := range items {
		fmt.Println(ii.Name)
		fmt.Println(ii.Hero)
		fmt.Println(steamCDN + ii.Image)
		//fmt.Println(ii.Description)
		fmt.Println(strings.Repeat("-", 55))
	}

	fmt.Println("total", len(items))
}

func getSource(steamID string) string {
	return filepath.Join(os.TempDir(), fmt.Sprintf(filenameFmt, steamID))
}

func dlCache(steamID string) error {
	url := fmt.Sprintf(steamInventoryAPI, steamID)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	filepath := getSource(steamID)
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
