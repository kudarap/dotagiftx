// Phantasm crawler code block can be paste on serverless function.
// Warning! Remember to replace package phantasm to package main

package phantasm

import (
	"bytes"
	"crypto/sha1"
	"encoding/gob"
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
	WebhookAuthHeader = "X-Require-Whisk-Auth"

	queryLimit   = 2000
	requestDelay = 1000 * time.Millisecond

	// ex. https://steamcommunity.com/inventory/76561198088587178/570/2
	steamURL            = "https://steamcommunity.com/inventory/%s/570/2?count=%d&start_assetid=%s"
	firstInventoryCount = 25
)

var (
	webhookURL string
	secret     string
)

func Main(args map[string]interface{}) map[string]interface{} {
	startsAt := time.Now()

	log.Println("starting phantasm...")
	if err := loadConfig(); err != nil {
		return resp(http.StatusInternalServerError, err)
	}

	id, ok := args["steam_id"]
	if !ok {
		return resp(http.StatusBadRequest, "missing steam_id")
	}
	steamID, ok := id.(string)
	if !ok {
		return resp(http.StatusBadRequest, "steam_id is not a string")
	}

	limit := queryLimit
	_, precheck := args["precheck"]
	if precheck {
		limit = firstInventoryCount
	}

	log.Println("starting requests...")
	var parts int
	var inventoryCount int
	var lastAssetID string
	var invent *inventory
	for {
		parts++
		log.Println("requesting part...", parts)
		next, status, err := get(steamID, limit, lastAssetID)
		if err != nil {
			return resp(status, err)
		}

		log.Println("merging inventories...")
		invent = merge(invent, next)

		lastAssetID = next.LastAssetID
		inventoryCount = next.TotalInventoryCount
		if next.MoreItems == 0 || precheck {
			break
		}
		time.Sleep(requestDelay)
	}

	var precheckHash string
	if precheck {
		log.Println("precheck - computing hash and skipping posting to webhook")
		precheckHash = invent.hash(steamID)
	} else {
		if err := post(steamID, invent); err != nil {
			return resp(http.StatusInternalServerError, err)
		}
	}

	log.Println("done!")
	return resp(http.StatusOK, map[string]interface{}{
		"precheck":         precheck,
		"precheck_hash":    precheckHash,
		"steam_id":         steamID,
		"query_limit":      limit,
		"request_delay_ms": requestDelay.Milliseconds(),
		"parts":            parts,
		"inventory_count":  inventoryCount,
		"elapsed_sec":      time.Since(startsAt).Seconds(),
		"webhook_url":      webhookURL,
		"last_asset_id":    lastAssetID,
	})
}

type inventory struct {
	Assets              []asset       `json:"assets"`
	Descriptions        []description `json:"descriptions"`
	TotalInventoryCount int           `json:"total_inventory_count"`

	LastAssetID string `json:"last_assetid"`
	MoreItems   int    `json:"more_items"`
	Rwgrsn      int    `json:"rwgrsn"`
	Success     int    `json:"success"`
}

func (i *inventory) hash(steamID string) string {
	var b bytes.Buffer
	gob.NewEncoder(&b).Encode(i)

	h := sha1.New()
	io.WriteString(h, steamID)
	io.WriteString(h, string(b.Bytes()))
	return fmt.Sprintf("%x", h.Sum(nil))
}

type asset struct {
	Amount     string `json:"amount"`
	AppID      int    `json:"appid"`
	AssetID    string `json:"assetid"`
	ClassID    string `json:"classid"`
	ContextID  string `json:"contextid"`
	InstanceID string `json:"instanceid"`
}

type description struct {
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

func merge(res ...*inventory) *inventory {
	var inventory inventory

	classIdx := map[string]struct{}{}
	for _, r := range res {
		if r == nil {
			continue
		}
		inventory.TotalInventoryCount = r.TotalInventoryCount
		inventory.Assets = append(inventory.Assets, r.Assets...)

		// map reduced descriptions across inventory.
		for _, d := range r.Descriptions {
			key := d.ClassID + "-" + d.InstanceID
			_, ok := classIdx[key]
			if !ok {
				classIdx[key] = struct{}{} // mark done
				inventory.Descriptions = append(inventory.Descriptions, d)
			}
		}
	}

	return &inventory
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
		if res.StatusCode == http.StatusForbidden {
			return nil, http.StatusForbidden, errors.New("private inventory")
		}
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
	req.Header.Set(WebhookAuthHeader, secret)

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
