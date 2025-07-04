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

	jsoniter "github.com/json-iterator/go"
	"github.com/kudarap/dotagiftx/steam"
)

const (
	defaultInventoryHashTTL = time.Hour * 2
	defaultRecrawlCD        = time.Minute * 10
	defaultCrawlerCD        = time.Minute

	maxWaitRetry = 5
)

var (
	errFileNotFound = fmt.Errorf("raw file not found")
	errFileWaiting  = fmt.Errorf("waiting for file")

	fastjson = jsoniter.ConfigFastest
)

type Service struct {
	id string

	config   Config
	cooldown cooldown
	logger   *slog.Logger

	recrawlCooldown  time.Duration
	crawlerCooldown  time.Duration
	inventoryHashTTL time.Duration

	electedCrawlerID int
}

func NewService(config Config, cd cooldown, logger *slog.Logger) *Service {
	config = config.setDefault()
	if err := os.MkdirAll(config.Path, 0777); err != nil {
		panic(err)
	}

	return &Service{
		id:               "phantasm",
		config:           config,
		cooldown:         cd,
		recrawlCooldown:  defaultRecrawlCD,
		crawlerCooldown:  defaultCrawlerCD,
		inventoryHashTTL: defaultInventoryHashTTL,
		logger:           logger.With("module", "phantasm"),
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

func (s *Service) InventoryAsset(ctx context.Context, steamID string) ([]steam.Asset, error) {
	raw, err := s.crawlWait(ctx, steamID)
	if err != nil {
		return nil, err
	}

	compat := raw.compat()
	return compat.ToAssets(), nil
}

func (s *Service) InventoryAssetWithProvider(ctx context.Context, steamID string) (string, []steam.Asset, error) {
	res, err := s.InventoryAsset(ctx, steamID)
	return s.id, res, err
}

// crawlWait retrieves the inventory local file when available and fetch it when missing.
func (s *Service) crawlWait(ctx context.Context, steamID string) (*inventory, error) {
	logger := s.logger.With("steam_id", steamID)
	logger.DebugContext(ctx, "fetch inventory and wait")

	localFile, err := s.localInventoryFile(ctx, steamID)
	if err != nil && !errors.Is(err, errFileNotFound) {
		return nil, err
	}
	if localFile != nil {
		logger.DebugContext(ctx, "local inventory ready")

		// when the hash still exists, the validity file age is still valid by inventoryHashExpr.
		hash, err := s.cooldown.InventoryHash(ctx, steamID)
		if err != nil {
			return nil, err
		}
		if hash != "" {
			// local file still good by hash and need to extend its validity.
			logger.DebugContext(ctx, "local inventory still fresh by hash signature",
				"hash", hash,
				"max_age", s.inventoryHashTTL,
				"extended_bt", s.inventoryHashTTL,
			)
			if err = s.cooldown.SetInventoryHash(ctx, steamID, hash, s.inventoryHashTTL); err != nil {
				return nil, err
			}
			return localFile, nil
		}

		// pre-check before fetching
		logger.DebugContext(ctx, "check remote inventory changes", "hash", hash)
		changed, err := s.remoteInventoryChanged(ctx, steamID)
		if err != nil {
			logger.Error("precheck remote inventory", "err", err)
			return nil, err
		}
		if !changed {
			logger.DebugContext(ctx, "remote inventory did not changed, falling back to local")
			// refresh inventory file age
			n := time.Now()
			if err = os.Chtimes(s.filePath(steamID), n, n); err != nil {
				return nil, err
			}
			return localFile, nil
		}

		logger.DebugContext(ctx, "local inventory requires re-fetch", "max_age", s.inventoryHashTTL)
	}

	// don't retry if it's not on waiting state.
	logger.DebugContext(ctx, "local file not found, crawling...")
	err = s.crawlRemoteInventory(ctx, steamID)
	if err != nil && !errors.Is(err, errFileWaiting) {
		return nil, err
	}
	// retry if its on waiting state.
	if errors.Is(err, errFileWaiting) {
		for i := range maxWaitRetry {
			wait := time.Duration(i*i) * time.Second
			time.Sleep(wait)
			logger.DebugContext(ctx, "reading local data", "attempt", i+1, "waiting", wait)
			localFile, err = s.localInventoryFile(ctx, steamID)
			if err != nil && !errors.Is(err, errFileNotFound) {
				return nil, err
			}
		}
	}
	// check raw inventory again but what error you have you need to go.
	localFile, err = s.localInventoryFile(ctx, steamID)
	if err != nil {
		return nil, err
	}

	// clear retry
	crawlerURL := s.config.Addrs[s.electedCrawlerID]
	crawlerID := extractCrawlerID(crawlerURL)
	if err = s.cooldown.SetRetryCooldown(ctx, crawlerID, steamID, 0); err != nil {
		return nil, err
	}

	return localFile, nil
}

func (s *Service) crawlRemoteInventory(ctx context.Context, steamID string) error {
	crawlerURL := s.config.Addrs[s.electedCrawlerID]
	crawlerID := extractCrawlerID(crawlerURL)
	logger := s.logger.With("steam_id", steamID, "crawler_id", crawlerID)

	// check if crawler is ready
	cd, err := s.cooldown.CrawlerCooldown(ctx, crawlerID)
	if err != nil {
		return err
	}
	if cd {
		return fmt.Errorf("crawler %s is on cooldown", crawlerID)
	}

	// check if there's existing requests
	cd, err = s.cooldown.RetryCooldown(ctx, crawlerID, steamID)
	if err != nil {
		return err
	}
	if cd {
		logger.DebugContext(ctx, "skipping crawling, please wait after", "recrawl_cd", s.recrawlCooldown)
		return errFileWaiting
	}
	if err = s.cooldown.SetRetryCooldown(ctx, crawlerID, steamID, s.recrawlCooldown); err != nil {
		return err
	}

	summary, err := s.sendCrawlRequest(ctx, crawlerURL, steamID, false)
	if err != nil {
		return err
	}
	logger.DebugContext(ctx,
		"fetch remote inventory",
		"count", summary.InventoryCount,
		"parts", summary.Parts,
		"query_limit", summary.QueryLimit,
		"request_delay_ms", summary.RequestDelayMs,
		"webhook_url", summary.WebhookURL,
	)
	return nil
}

func (s *Service) remoteInventoryChanged(ctx context.Context, steamID string) (bool, error) {
	crawlerURL := s.config.Addrs[s.electedCrawlerID]
	crawlerID := extractCrawlerID(crawlerURL)
	logger := s.logger.With("steam_id", steamID, "crawler_id", crawlerID)

	result, err := s.sendCrawlRequest(ctx, crawlerURL, steamID, true)
	if err != nil {
		return false, err
	}

	logger.DebugContext(ctx, "precheck remote inventory", "hash", result.PrecheckHash)
	currentHash, err := s.cooldown.InventoryHash(ctx, steamID)
	if err != nil {
		return false, err
	}

	logger.DebugContext(ctx, "comparing hashes", "current", currentHash, "new", result.PrecheckHash)
	if err = s.cooldown.SetInventoryHash(ctx, steamID, result.PrecheckHash, s.inventoryHashTTL); err != nil {
		return true, err
	}
	return result.PrecheckHash != currentHash, nil
}

func (s *Service) localInventoryFile(ctx context.Context, steamID string) (*inventory, error) {
	file, err := os.ReadFile(s.filePath(steamID))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errFileNotFound
		}
		return nil, fmt.Errorf("open file: %s", err)
	}

	var inv inventory
	if err = fastjson.Unmarshal(file, &inv); err != nil {
		return nil, fmt.Errorf("unmarshal: %s", err)
	}
	return &inv, nil
}

func (s *Service) electNewCrawler(ctx context.Context) string {
	crawler := extractCrawlerID(s.config.Addrs[s.electedCrawlerID])
	cd, err := s.cooldown.CrawlerCooldown(ctx, crawler)
	if err != nil {
		s.logger.ErrorContext(ctx, "crawler cooldown", "crawler", crawler, "err", err)
	}
	if !cd {
		if err = s.cooldown.SetCrawlerCooldown(ctx, crawler, s.crawlerCooldown); err != nil {
			s.logger.ErrorContext(ctx, "crawler cooldown", "crawler", crawler, "err", err)
		}
	}

	s.electedCrawlerID++
	if s.electedCrawlerID >= len(s.config.Addrs) {
		s.electedCrawlerID = 0
	}
	return extractCrawlerID(s.config.Addrs[s.electedCrawlerID])
}

func (s *Service) sendCrawlRequest(ctx context.Context, crawlerURL, steamID string, precheck bool) (*CrawlSummary, error) {
	url := fmt.Sprintf("%s?steam_id=%s", crawlerURL, steamID)
	if precheck {
		url = fmt.Sprintf("%s&precheck", url)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = s.setRequestHeaders()

	var summary CrawlSummary
	statusCode, err := sendRequest(req, &summary)
	if err != nil {
		if statusCode == http.StatusForbidden {
			return nil, steam.ErrInventoryPrivate
		}

		// only elect new crawler when not found and too much request
		if statusCode == http.StatusNotFound || statusCode == http.StatusTooManyRequests {
			elected := s.electNewCrawler(ctx)
			s.logger.InfoContext(ctx, "current crawler unavailable, new crawler elected",
				"old", extractCrawlerID(crawlerURL),
				"new", elected,
			)
		}
		return nil, err
	}
	return &summary, nil
}

func (s *Service) setRequestHeaders() http.Header {
	h := http.Header{}
	h.Add("Content-Type", "application/json")
	h.Add("X-Require-Whisk-Auth", s.config.Secret)
	return h
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

type cooldown interface {
	RetryCooldown(ctx context.Context, crawlID, steamID string) (bool, error)
	SetRetryCooldown(ctx context.Context, crawlID, steamID string, ttl time.Duration) error

	CrawlerCooldown(ctx context.Context, crawlID string) (bool, error)
	SetCrawlerCooldown(ctx context.Context, crawlID string, ttl time.Duration) error

	InventoryHash(ctx context.Context, steamID string) (hash string, error error)
	SetInventoryHash(ctx context.Context, steamID, hash string, ttl time.Duration) error
}

func extractCrawlerID(addr string) string {
	ss := strings.Split(addr, "/")
	return ss[len(ss)-1]
}
