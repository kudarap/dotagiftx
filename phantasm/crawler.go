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
	var invs []*inventory
	for {
		parts++
		log.Println("requesting part...", parts)
		res, status, err := get(steamID, queryLimit, startAssetID)
		if err != nil {
			return resp(status, err)
		}
		startAssetID = res.LastAssetid
		invs = append(invs, res)
		if res.MoreItems == 0 {
			inventoryCount = res.TotalInventoryCount
			break
		}
		time.Sleep(requestDelay)
	}

	// combine data here
	log.Println("combining inventories...")
	inv := merge(invs...)
	if err := post(steamID, inv); err != nil {
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

type inventory struct {
	Assets              []asset       `json:"assets"`
	Descriptions        []description `json:"descriptions"`
	TotalInventoryCount int           `json:"total_inventory_count"`

	LastAssetid string `json:"last_assetid"`
	MoreItems   int    `json:"more_items"`
	Rwgrsn      int    `json:"rwgrsn"`
	Success     int    `json:"success"`
}

type asset struct {
	Amount     string `json:"amount"`
	Appid      int    `json:"appid"`
	Assetid    string `json:"assetid"`
	Classid    string `json:"classid"`
	Contextid  string `json:"contextid"`
	Instanceid string `json:"instanceid"`
}

type description struct {
	Appid           int    `json:"appid"`
	BackgroundColor string `json:"background_color"`
	Classid         string `json:"classid"`
	Commodity       int    `json:"commodity"`
	Currency        int    `json:"currency"`
	Descriptions    []struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"descriptions"`
	IconURL                     string `json:"icon_url"`
	IconURLLarge                string `json:"icon_url_large"`
	Instanceid                  string `json:"instanceid"`
	MarketHashName              string `json:"market_hash_name"`
	MarketMarketableRestriction int    `json:"market_marketable_restriction"`
	MarketName                  string `json:"market_name"`
	MarketTradableRestriction   int    `json:"market_tradable_restriction"`
	Marketable                  int    `json:"marketable"`
	Name                        string `json:"name"`
	NameColor                   string `json:"name_color"`
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

func merge(res ...*inventory) *inventory {
	var inv inventory

	classIdx := map[string]struct{}{}
	for _, r := range res {
		if res == nil {
			continue
		}
		inv.TotalInventoryCount = r.TotalInventoryCount
		inv.Assets = append(inv.Assets, r.Assets...)

		// map reduced descriptions across inventory.
		for _, d := range r.Descriptions {
			key := d.Classid + "-" + d.Instanceid
			_, ok := classIdx[key]
			if !ok {
				classIdx[key] = struct{}{} // mark done
				inv.Descriptions = append(inv.Descriptions, d)
			}
		}
	}

	return &inv
}

func get(steamID string, count int, lastAssetID string) (*inventory, int, error) {
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

	var inv inventory
	if err = json.Unmarshal(body, &inv); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return &inv, http.StatusOK, nil
}

func post(steamID string, inv *inventory) error {
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
