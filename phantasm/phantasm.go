// Phantasm Crawler
//
// "Drawing on his battles fought across many worlds and many times, phantasms of the Chaos Knight rise up to quell all
// who oppose him"
//
// "Summons several phantasmal copies of the Chaos Knight from alternate dimensions. The phantasms are illusions that
// deal 100% damage, but take 350% damage."
//
// Phantasm crawls Inventory for item and delivery tracking. Hopefully, by summoning multiple instances of the crawler
// will provide better steam Inventory raw data.
//
// crawler.go
//	- script is intended for serverless functions to work around with ip rate limits during peak usage.
// 	- publishes raw Inventory data to target webhook url.

package phantasm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/djherbis/times"
	jsoniter "github.com/json-iterator/go"
	"github.com/kudarap/dotagiftx/steam"
)

var (
	errFileNotFound = fmt.Errorf("raw file not found")
	errFileWaiting  = fmt.Errorf("waiting for file")

	fastjson = jsoniter.ConfigFastest
)

type Config struct {
	Addrs      []string
	WebhookURL string `envconfig:"WEBHOOK_URL"`
	Secret     string
	Path       string
}

type Service struct {
	config      Config
	cachePrefix string
	cacheTTL    time.Duration
	logger      *slog.Logger

	electedCrawlerID int
	retryAfter       map[string]time.Time
}

func (s *Service) InventoryAsset(steamID string) ([]steam.Asset, error) {
	ctx := context.Background()
	raw, err := s.autoRetry(ctx, steamID)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}

	compat := raw.compat()
	return compat.ToAssets(), nil
}

func (s *Service) SaveInventory(ctx context.Context, steamID string, body io.ReadCloser) error {
	file, err := os.Create(s.filePath(steamID))
	if err != nil {
		return fmt.Errorf("open file: %s", err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			s.logger.Error("close file", "error", err.Error())
		}
		if err = body.Close(); err != nil {
			s.logger.Error("close body", "error", err.Error())
		}
	}()
	if _, err = io.Copy(file, body); err != nil {
		return fmt.Errorf("copy: %s", err)
	}

	return nil
}

func (s *Service) autoRetry(ctx context.Context, steamID string) (*Inventory, error) {
	raw, err := s.rawInventory(ctx, steamID)
	if err != nil && !errors.Is(err, errFileNotFound) {
		return nil, err
	}

	err = s.crawlInventory(ctx, steamID)
	if err != nil && !errors.Is(err, errFileWaiting) {
		return nil, err
	}

	for i := range 5 {
		time.Sleep(time.Duration(i+5) * time.Second)
		s.logger.Info("retrying steam", "attempt", i)
		raw, err = s.rawInventory(ctx, steamID)
		if err != nil && !errors.Is(err, errFileNotFound) {
			return nil, err
		}
	}

	// re-fetch day old file
	t, err := times.Stat(s.filePath(steamID))
	if err != nil {
		return nil, err
	}
	log.Println(t.ModTime())

	return raw, nil
}

func crawlerName(addr string) string {
	ss := strings.Split(addr, "/")
	return ss[len(ss)-1]
}

func (s *Service) crawlInventory(ctx context.Context, steamID string) error {
	crawlerURL := s.config.Addrs[s.electedCrawlerID]
	crawlerID := crawlerName(crawlerURL)

	lastReq, ok := s.retryAfter[crawlerID]
	if ok && time.Since(lastReq) < s.cacheTTL {
		return errFileWaiting
	}
	s.retryAfter[crawlerID] = time.Now()

	s.logger.InfoContext(ctx, "elected crawler", "crawler", crawlerID, "steamID", steamID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, crawlerURL+"?steam_id="+steamID, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Require-Whisk-Auth", s.config.Secret)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.ErrorContext(ctx, "fetch raw inventory", "steam_id", steamID, "error", err)
		return err
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			s.logger.ErrorContext(ctx, "close body", "error", err.Error())
		}
	}()
	if res.StatusCode > 299 {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("%d - %s", res.StatusCode, body)
	}

	data := struct {
		ElapsedSec     float64 `json:"elapsed_sec"`
		InventoryCount int     `json:"inventory_count"`
		Parts          int     `json:"parts"`
		QueryLimit     int     `json:"query_limit"`
		RequestDelayMs int     `json:"request_delay_ms"`
		SteamID        string  `json:"steam_id"`
		WebhookURL     string  `json:"webhook_url"`
	}{}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, &data); err != nil {
		return err
	}

	s.logger.Info("fetch raw inventory",
		"steam_id", steamID,
		"count", data.InventoryCount,
		"parts", data.Parts,
		"query_limit", data.QueryLimit,
		"request_delay_ms", data.RequestDelayMs,
		"steam_id", steamID,
		"webhook_url", data.WebhookURL,
	)
	return nil
}

func (s *Service) rawInventory(ctx context.Context, steamID string) (*Inventory, error) {
	file, err := os.ReadFile(s.filePath(steamID))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errFileNotFound
		}
		return nil, fmt.Errorf("open file: %s", err)
	}

	var inventory Inventory
	if err = fastjson.Unmarshal(file, &inventory); err != nil {
		return nil, fmt.Errorf("unmarshal: %s", err)
	}
	return &inventory, nil
}

func (s *Service) filePath(steamID string) string {
	return filepath.Join(s.config.Path, fmt.Sprintf("%s.json", steamID))
}

func NewService(config Config, logger *slog.Logger) *Service {
	if err := os.MkdirAll(config.Path, 0777); err != nil {
		panic(err)
	}

	for i, addr := range config.Addrs {
		fmt.Println(i, addr)
	}

	return &Service{
		config:      config,
		cachePrefix: "phantasm",
		cacheTTL:    time.Hour,
		logger:      logger,
		retryAfter:  map[string]time.Time{},
	}
}

func (i *Inventory) compat() steam.AllInventory {
	assets := make([]steam.RawInventoryAsset, len(i.Assets))
	for k, v := range i.Assets {
		assets[k] = v.compat()
	}

	descs := make(map[string]steam.RawInventoryDesc, len(i.Descriptions))
	for _, v := range i.Descriptions {
		descs[fmt.Sprintf("%s_%s", v.ClassID, v.InstanceID)] = v.compat()
	}

	return steam.AllInventory{AllInvs: assets, AllDescs: descs}
}

func (a *Asset) compat() steam.RawInventoryAsset {
	return steam.RawInventoryAsset{
		ID:         a.AssetID,
		AssetID:    a.AssetID,
		ClassID:    a.ClassID,
		InstanceID: a.InstanceID,
	}
}

func (d *Description) compat() steam.RawInventoryDesc {
	attrs := make(steam.RawInventoryItemDetails, len(d.Descriptions))
	for i, v := range d.Descriptions {
		attrs[i].Value = v.Value
	}

	return steam.RawInventoryDesc{
		ClassID:      d.ClassID,
		InstanceID:   d.InstanceID,
		Name:         d.Name,
		Image:        d.IconURLLarge,
		Type:         d.Type,
		Descriptions: attrs,
	}
}
