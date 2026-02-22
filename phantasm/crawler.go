// Phantasm crawler code block can be paste on serverless function.
// Warning! Remember to replace package phantasm to package main

package phantasm

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"
)

const (
	WebhookAuthHeader = "X-Require-Whisk-Auth"

	precheckLimit = 25
	queryLimit    = 2000
	requestDelay  = 1000 * time.Millisecond

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

	ctx := context.Background()
	now := time.Now()
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
		limit = precheckLimit
	}

	log.Println("starting requests...")
	var parts int
	var inventoryCount int
	var lastAssetID string
	var invent *inventory
	for {
		parts++
		log.Println("requesting part...", parts)
		next, status, err := get(ctx, steamID, limit, lastAssetID)
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
		h, err := invent.hash(steamID)
		if err != nil {
			return resp(http.StatusInternalServerError, err)
		}
		precheckHash = h
	} else {
		if err := post(ctx, steamID, invent); err != nil {
			log.Println("posting failed:", err)
			return resp(http.StatusInternalServerError, err)
		}
	}

	log.Println("done!")
	var summary CrawlSummary
	summary.SteamID = steamID
	summary.QueryLimit = limit
	summary.RequestDelayMs = int(requestDelay.Milliseconds())
	summary.Parts = parts
	summary.InventoryCount = inventoryCount
	summary.ElapsedSec = time.Since(now).Seconds()
	summary.WebhookURL = webhookURL
	summary.Precheck = precheck
	summary.PrecheckHash = precheckHash
	summary.LastAssetID = lastAssetID
	return resp(http.StatusOK, structToMap(summary))
}

type CrawlSummary struct {
	ElapsedSec     float64 `json:"elapsed_sec"`
	InventoryCount int     `json:"inventory_count"`
	Parts          int     `json:"parts"`
	QueryLimit     int     `json:"query_limit"`
	RequestDelayMs int     `json:"request_delay_ms"`
	SteamID        string  `json:"steam_id"`
	WebhookURL     string  `json:"webhook_url"`
	Precheck       bool    `json:"precheck"`
	PrecheckHash   string  `json:"precheck_hash"`
	LastAssetID    string  `json:"last_asset_id"`
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

func (i *inventory) hash(steamID string) (string, error) {
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(i); err != nil {
		return "", err
	}

	h := sha1.New()
	h.Write([]byte(steamID + b.String()))
	return hex.EncodeToString(h.Sum(nil)), nil
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
	Descriptions                []descriptionAttrs `json:"descriptions"`
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

type descriptionAttrs struct {
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

func get(ctx context.Context, steamID string, count int, lastAssetID string) (i *inventory, statusCode int, err error) {
	url := fmt.Sprintf(steamURL, steamID, count, lastAssetID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var inv inventory
	statusCode, err = sendRequest(req, &inv)
	if err != nil {
		return nil, statusCode, err
	}
	return &inv, statusCode, nil
}

func post(ctx context.Context, steamID string, data any) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	url := strings.TrimRight(webhookURL, "/") + "/" + steamID
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(WebhookAuthHeader, secret)

	if _, err = sendRequest(req, nil); err != nil {
		return err
	}
	return nil
}

func sendRequest(req *http.Request, out any) (statusCode int, err error) {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	res.Body.Close()
	if res.StatusCode > 299 {
		return res.StatusCode, errors.New(http.StatusText(res.StatusCode))
	}

	if out != nil {
		if err = json.Unmarshal(body, &out); err != nil {
			return 0, err
		}
	}
	return res.StatusCode, nil
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

func structToMap(data interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	if data == nil {
		return res
	}
	v := reflect.TypeOf(data)
	reflectValue := reflect.ValueOf(data)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")
		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = structToMap(field)
			} else {
				res[tag] = field
			}
		}
	}
	return res
}
