// Phantasm Crawler
//
// "Drawing on his battles fought across many worlds and many times, phantasms of the Chaos Knight rise up to quell all
// who oppose him"
//
// "Summons several phantasmal copies of the Chaos Knight from alternate dimensions. The phantasms are illusions that
// deal 100% damage, but take 350% damage."
//
// Phantasm crawls inventory for item and delivery tracking. Hopefully, by summoning multiple instances of the crawler
// will provide better steam inventory raw data.
//
// crawler.go
//	- script is intended for serverless functions to work around with ip rate limits during peak usage.
// 	- publishes raw inventory data to target webhook url.

package phantasm

import (
	"context"
	"errors"
	"fmt"
	"io"
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

const (
	defaultMaxAge        = time.Hour * 12
	defaultRetryCooldown = time.Hour
	defaultCrawlCooldown = time.Minute
)

var (
	errFileNotFound = fmt.Errorf("raw file not found")
	errFileWaiting  = fmt.Errorf("waiting for file")

	fastjson = jsoniter.ConfigFastest
)

type Service struct {
	config        Config
	cachePrefix   string
	retryCooldown time.Duration
	maxAge        time.Duration
	logger        *slog.Logger

	electedCrawlerID int
	retryAfter       map[string]time.Time
	crawlerCoolAfter map[string]time.Time
	crawlerCooldown  time.Duration
}

func NewService(config Config, logger *slog.Logger) *Service {
	config = config.setDefault()
	if err := os.MkdirAll(config.Path, 0777); err != nil {
		panic(err)
	}

	return &Service{
		config:           config,
		cachePrefix:      "phantasm",
		maxAge:           defaultMaxAge,
		logger:           logger,
		retryAfter:       map[string]time.Time{},
		retryCooldown:    defaultRetryCooldown,
		crawlerCoolAfter: map[string]time.Time{},
		crawlerCooldown:  defaultCrawlCooldown,
	}
}

func (s *Service) SaveInventory(ctx context.Context, steamID, secret string, body io.ReadCloser) error {
	// ensure that the filename has no path separators or parent directory references
	if steamID == "" || strings.Contains(steamID, "/") || strings.Contains(steamID, "\\") ||
		strings.Contains(steamID, "..") {
		return errors.New("invalid steam id")
	}
	if secret != s.config.Secret {
		return fmt.Errorf("invalid secret")
	}

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

func (s *Service) InventoryAsset(steamID string) ([]steam.Asset, error) {
	ctx := context.Background()
	raw, err := s.autoRetry(ctx, steamID)
	if err != nil {
		return nil, err
	}

	compat := raw.compat()
	return compat.ToAssets(), nil
}

func (s *Service) autoRetry(ctx context.Context, steamID string) (*inventory, error) {
	invent, err := s.rawInventory(ctx, steamID)
	if err != nil && !errors.Is(err, errFileNotFound) {
		return nil, err
	}
	if invent != nil {
		// re-fetch day old file
		t, err := times.Stat(s.filePath(steamID))
		if err != nil {
			return nil, err
		}
		age := time.Since(t.ModTime())
		if age < s.maxAge {
			return invent, nil
		}
		s.logger.Info("max age reached, recrawl", "steamid", steamID, "age", age, "max-age", s.maxAge)
	}

	err = s.crawlInventory(ctx, steamID)
	// don't retry if it's not on waiting state.
	if err != nil && !errors.Is(err, errFileWaiting) {
		return nil, err
	}
	// retry if its on waiting state.
	if errors.Is(err, errFileWaiting) {
		for i := range 5 {
			wait := time.Duration(i+1) * time.Second
			time.Sleep(wait)
			s.logger.Info("retrying steam", "attempt", i+1, "steamid", steamID, "waiting", wait)
			invent, err = s.rawInventory(ctx, steamID)
			if err != nil && !errors.Is(err, errFileNotFound) {
				return nil, err
			}
		}
	}
	// check raw inventory again but what error you have you need to go.
	invent, err = s.rawInventory(ctx, steamID)
	if err != nil {
		return nil, err
	}

	// clear retry
	crawlerURL := s.config.Addrs[s.electedCrawlerID]
	crawlerID := crawlerName(crawlerURL)
	retryID := crawlerID + "-" + steamID
	delete(s.retryAfter, retryID)
	return invent, nil
}

func (s *Service) crawlInventory(ctx context.Context, steamID string) error {
	timeNow := time.Now()
	crawlerURL := s.config.Addrs[s.electedCrawlerID]
	crawlerID := crawlerName(crawlerURL)
	s.logger.InfoContext(ctx, "elected crawler", "crawler", crawlerID, "steamID", steamID)

	// check if crawler is ready
	cd, ok := s.crawlerCoolAfter[crawlerID]
	if ok && cd.After(timeNow) {
		return fmt.Errorf("crawler %s is on all cooldown", crawlerID)
	}

	// check if there's existing requests
	retryID := crawlerID + "-" + steamID
	lastReq, ok := s.retryAfter[retryID]
	if ok && lastReq.After(timeNow) {
		s.logger.InfoContext(ctx, "skipping crawling, please wait after",
			"last_req", lastReq, "ttl", s.retryCooldown,
		)
		return errFileWaiting
	}
	s.retryAfter[retryID] = timeNow.Add(s.retryCooldown)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, crawlerURL+"?steam_id="+steamID, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Require-Whisk-Auth", s.config.Secret)
	var summary CrawlSummary
	statusCode, err := sendRequest(req, &summary)
	if err != nil {
		// only elect new crawler when not found and too much request
		if statusCode == http.StatusNotFound || statusCode == http.StatusTooManyRequests {
			s.electNewCrawler()
		}
		if statusCode == http.StatusForbidden {
			return steam.ErrInventoryPrivate
		}
		return err
	}
	s.logger.Info("fetch raw inventory",
		"steam_id", summary.SteamID,
		"count", summary.InventoryCount,
		"parts", summary.Parts,
		"query_limit", summary.QueryLimit,
		"request_delay_ms", summary.RequestDelayMs,
		"webhook_url", summary.WebhookURL,
	)
	return nil
}

func (s *Service) rawInventory(ctx context.Context, steamID string) (*inventory, error) {
	file, err := os.ReadFile(s.filePath(steamID))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errFileNotFound
		}
		return nil, fmt.Errorf("open file: %s", err)
	}

	var inventory inventory
	if err = fastjson.Unmarshal(file, &inventory); err != nil {
		return nil, fmt.Errorf("unmarshal: %s", err)
	}
	return &inventory, nil
}

func (s *Service) electNewCrawler() {
	current := s.config.Addrs[s.electedCrawlerID]
	id := crawlerName(current)
	if _, ok := s.crawlerCoolAfter[id]; !ok {
		s.crawlerCoolAfter[id] = time.Now().Add(s.crawlerCooldown)
		s.electedCrawlerID++
		if s.electedCrawlerID >= len(s.config.Addrs) {
			s.electedCrawlerID = 0
		}
	}
}

func (s *Service) filePath(steamID string) string {
	return filepath.Join(s.config.Path, fmt.Sprintf("%s.json", steamID))
}

func (i *inventory) compat() steam.AllInventory {
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

func (a *asset) compat() steam.RawInventoryAsset {
	return steam.RawInventoryAsset{
		ID:         a.AssetID,
		AssetID:    a.AssetID,
		ClassID:    a.ClassID,
		InstanceID: a.InstanceID,
	}
}

func (d *description) compat() steam.RawInventoryDesc {
	attrs := make(steam.RawInventoryItemDetails, len(d.Descriptions))
	for i, v := range d.Descriptions {
		attrs[i].Value = strings.TrimPrefix(v.Value, "\n")
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

func crawlerName(addr string) string {
	ss := strings.Split(addr, "/")
	return ss[len(ss)-1]
}
