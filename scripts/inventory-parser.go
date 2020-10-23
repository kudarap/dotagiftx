package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigFastest

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

func main() {
	fmt.Println("hello")

	// read file
	data, err := ioutil.ReadFile("./inventory.json")
	if err != nil {
		fmt.Println("could not read file:", err)
		return
	}

	var inv Inventory
	if err := json.Unmarshal(data, &inv); err != nil {
		fmt.Println("could not parse json:", err)
	}
	fmt.Println("parsed", len(inv.Items))

	const filter = "International 2019"
	items := map[string]Item{}
	for _, ii := range inv.Items {
		//if ii.Type != "Mythical Bundle" {
		if !strings.Contains(ii.Type, "Mythical Bundle") {
			continue
		}

		desc, hero := ii.stringifyDesc()
		if !strings.Contains(strings.ToLower(desc), strings.ToLower(filter)) {
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
		fmt.Println(ii.Image)
		//fmt.Println(ii.Description)
		fmt.Println(strings.Repeat("-", 55))
	}

	fmt.Println("total", len(items))
}
