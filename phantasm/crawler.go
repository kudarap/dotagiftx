// Phantasm crawler code block can be paste on serverless function.
// Warning! Remember to replace package phantasm to package main

package phantasm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	queryLimit   = 2000
	requestDelay = 500 * time.Millisecond

	webhookAuthHeader = "X-Require-Whisk-Auth"

	// ex. https://steamcommunity.com/inventory/76561198088587178/570/2
	steamURL = "https://steamcommunity.com/inventory/%s/570/2?count=%d&start_assetid=%s"
)

var (
	webhookURL string
	secret     string
)

func Main(args map[string]interface{}) map[string]interface{} {
	log.Println("starting phantasm...")
	if err := loadConfig(); err != nil {
		return resp(http.StatusInternalServerError, err)
	}

	now := time.Now()
	id, ok := args["steam_id"]
	if !ok {
		return resp(http.StatusBadRequest, "missing steam_id")
	}
	steamID, ok := id.(string)
	if !ok {
		return resp(http.StatusBadRequest, "steam_id is not a string")
	}

	log.Println("starting requests...")
	var parts int
	var inventoryCount int
	var startAssetID string
	var inventory *Inventory
	for {
		parts++
		log.Println("requesting part...", parts)
		next, status, err := get(steamID, queryLimit, startAssetID)
		if err != nil {
			return resp(status, err)
		}

		log.Println("merging inventories...")
		inventory = merge(inventory, next)

		startAssetID = next.LastAssetID
		if next.MoreItems == 0 {
			inventoryCount = next.TotalInventoryCount
			break
		}
		time.Sleep(requestDelay)
	}

	if err := post(steamID, inventory); err != nil {
		return resp(http.StatusInternalServerError, err)
	}

	log.Println("done!")
	return resp(http.StatusOK, map[string]interface{}{
		"steam_id":         steamID,
		"query_limit":      queryLimit,
		"request_delay_ms": requestDelay.Milliseconds(),
		"parts":            parts,
		"inventory_count":  inventoryCount,
		"elapsed_sec":      time.Since(now).Seconds(),
		"webhook_url":      webhookURL,
	})
}

type Inventory struct {
	Assets              []Asset       `json:"assets"`
	Descriptions        []Description `json:"descriptions"`
	TotalInventoryCount int           `json:"total_inventory_count"`

	LastAssetID string `json:"last_assetid"`
	MoreItems   int    `json:"more_items"`
	Rwgrsn      int    `json:"rwgrsn"`
	Success     int    `json:"success"`
}

type Asset struct {
	Amount     string `json:"amount"`
	AppID      int    `json:"appid"`
	AssetID    string `json:"assetid"`
	ClassID    string `json:"classid"`
	ContextID  string `json:"contextid"`
	InstanceID string `json:"instanceid"`
}

type Description struct {
	AppID                       int                `json:"appid"`
	BackgroundColor             string             `json:"background_color"`
	ClassID                     string             `json:"classid"`
	Commodity                   int                `json:"commodity"`
	Currency                    int                `json:"currency"`
	Descriptions                []DescriptionAttrs `json:"descriptions"`
	IconURL                     string             `json:"icon_url"`
	IconURLLarge                string             `json:"icon_url_large"`
	InstanceID                  string             `json:"instanceid"`
	MarketHashName              string             `json:"market_hash_name"`
	MarketMarketableRestriction int                `json:"market_marketable_restriction"`
	MarketName                  string             `json:"market_name"`
	MarketTradableRestriction   int                `json:"market_tradable_restriction"`
	Marketable                  int                `json:"marketable"`
	Name                        string             `json:"name"`
	NameColor                   string             `json:"name_color"`
	Tags                        []struct {
		Category              string `json:"category"`
		Color                 string `json:"color,omitempty"`
		InternalName          string `json:"internal_name"`
		LocalizedCategoryName string `json:"localized_category_name"`
		LocalizedTagName      string `json:"localized_tag_name"`
	} `json:"tags"`
	Tradable int    `json:"tradable"`
	Type     string `json:"type"`
}

type DescriptionAttrs struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

func merge(res ...*Inventory) *Inventory {
	var inv Inventory

	classIdx := map[string]struct{}{}
	for _, r := range res {
		if r == nil {
			continue
		}
		inv.TotalInventoryCount = r.TotalInventoryCount
		inv.Assets = append(inv.Assets, r.Assets...)

		// map reduced descriptions across Inventory.
		for _, d := range r.Descriptions {
			key := d.ClassID + "-" + d.InstanceID
			_, ok := classIdx[key]
			if !ok {
				classIdx[key] = struct{}{} // mark done
				inv.Descriptions = append(inv.Descriptions, d)
			}
		}
	}

	return &inv
}

func get(steamID string, count int, lastAssetID string) (*Inventory, int, error) {
	res, err := http.Get(fmt.Sprintf(steamURL, steamID, count, lastAssetID))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	res.Body.Close()

	if res.StatusCode > 299 {
		return nil, res.StatusCode, fmt.Errorf("%d - %s", res.StatusCode, body)
	}

	var inv Inventory
	if err = json.Unmarshal(body, &inv); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return &inv, http.StatusOK, nil
}

func post(steamID string, inv *Inventory) error {
	b, err := json.Marshal(inv)
	if err != nil {
		return err
	}

	url := strings.TrimRight(webhookURL, "/") + "/" + steamID
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(webhookAuthHeader, secret)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	res.Body.Close()
	if res.StatusCode > 299 {
		return fmt.Errorf("%d - %s", res.StatusCode, body)
	}
	return nil
}

func resp(status int, body interface{}) map[string]interface{} {
	return map[string]interface{}{
		"statusCode": status,
		"body":       body,
	}
}

func loadConfig() error {
	webhookURL = os.Getenv("DG_PHANTASM_WEBHOOK_URL")
	secret = os.Getenv("DG_PHANTASM_SECRET")
	if webhookURL == "" || secret == "" {
		return errors.New("webhookURL and secret required")
	}
	return nil
}
