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
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/kudarap/dotagiftx/steam"
)

var fastjson = jsoniter.ConfigFastest

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

	return steam.AllInventory{assets, descs}
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

func (s *Service) InventoryAsset(steamID string) ([]steam.Asset, error) {
	ctx := context.Background()
	raw, err := s.GetInventory(ctx, steamID)
	if err != nil {
		return nil, err
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

func (s *Service) GetInventory(ctx context.Context, steamID string) (*Inventory, error) {
	file, err := os.ReadFile(s.filePath(steamID))
	if err != nil {
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

	return &Service{
		config:      config,
		cachePrefix: "phantasm",
		cacheTTL:    time.Hour,
		logger:      logger,
	}
}
