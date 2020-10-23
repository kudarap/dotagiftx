package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigFastest

type Item struct {
	ClassID string `json:"classid"`
	Name    string `json:"name"`
	Image   string `json:"icon_url_large"`
	Type    string `json:"type"`
	Descs   []struct {
		Value string `json:"value"`
	} `json:"descriptionsX"`
	Des1        interface{} `json:"descriptions"`
	Description string      `json:"description"`
}

func (i Item) stringifyDesc() string {
	var s []string
	for _, dd := range i.Descs {
		s = append(s, dd.Value)
	}

	return strings.Join(s, " ")
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

		desc := fmt.Sprintf("%s", ii.Des1)
		if !strings.Contains(strings.ToLower(desc), strings.ToLower(filter)) {
			continue
		}

		items[ii.Name] = Item{
			ClassID:     ii.ClassID,
			Name:        ii.Name,
			Type:        ii.Type,
			Image:       ii.Image,
			Description: desc,
		}
	}

	for _, ii := range items {
		fmt.Println(ii.Name)
		fmt.Println(ii.Image)
		fmt.Println(ii.Description)
		fmt.Println(strings.Repeat("-", 55))
	}

	fmt.Println("total", len(items))
}
