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
	defaultMaxAge           = time.Hour * 12
	defaultRetryCooldown    = time.Minute * 5
	defaultCrawlCooldown    = time.Minute
	defaultInventoryHashTTL = time.Hour * 24 * 30
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

	retryCooldown   time.Duration
	maxAge          time.Duration
	crawlerCooldown time.Duration
	hashTTL         time.Duration

	electedCrawlerID int
}

func NewService(config Config, cd cooldown, logger *slog.Logger) *Service {
	config = config.setDefault()
	if err := os.MkdirAll(config.Path, 0777); err != nil {
		panic(err)
	}

	return &Service{
		id:              "phantasm",
		config:          config,
		cooldown:        cd,
		maxAge:          defaultMaxAge,
		retryCooldown:   defaultRetryCooldown,
		crawlerCooldown: defaultCrawlCooldown,
		hashTTL:         defaultInventoryHashTTL,
		logger:          logger.With("module", "phantasm"),
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
	raw, err := s.autoRetry(ctx, steamID)
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

		// pre-check before fetching
		s.logger.Info("checking inventory changed base on hash", "steamid", steamID)
		changed, err := s.hasInventoryChanged(ctx, steamID)
		if err != nil {
			return nil, err
		}
		if !changed {
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
			s.logger.Info("reading local data", "attempt", i+1, "steamid", steamID, "waiting", wait)
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
	if err = s.cooldown.SetRetryCooldown(ctx, crawlerID, steamID, 0); err != nil {
		return nil, err
	}

	return invent, nil
}

func (s *Service) crawlInventory(ctx context.Context, steamID string) error {
	crawlerURL := s.config.Addrs[s.electedCrawlerID]
	crawlerID := crawlerName(crawlerURL)
	s.logger.InfoContext(ctx, "elected crawler", "crawler", crawlerID, "steamID", steamID)

	// check if crawler is ready
	cd, err := s.cooldown.CrawlerCooldown(ctx, crawlerID)
	if err != nil {
		return err
	}
	if cd {
		return fmt.Errorf("crawler %s is on all cooldown", crawlerID)
	}

	// check if there's existing requests
	cd, err = s.cooldown.RetryCooldown(ctx, crawlerID, steamID)
	if err != nil {
		return err
	}
	if cd {
		s.logger.InfoContext(ctx, "skipping crawling, please wait after", "ttl", s.retryCooldown)
		return errFileWaiting
	}
	if err = s.cooldown.SetRetryCooldown(ctx, crawlerID, steamID, s.retryCooldown); err != nil {
		return err
	}

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
			s.electNewCrawler(ctx)
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

func (s *Service) hasInventoryChanged(ctx context.Context, steamID string) (bool, error) {
	return false, nil
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

func (s *Service) electNewCrawler(ctx context.Context) {
	crawler := crawlerName(s.config.Addrs[s.electedCrawlerID])
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

func crawlerName(addr string) string {
	ss := strings.Split(addr, "/")
	return ss[len(ss)-1]
}
