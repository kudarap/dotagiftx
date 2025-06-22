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
	"time"

	jsoniter "github.com/json-iterator/go"
	localcache "github.com/kudarap/dotagiftx/cache"
	"github.com/kudarap/dotagiftx/steam"
)

var fastjson = jsoniter.ConfigFastest

type Config struct {
	Addrs      []string
	WebhookURL string `envconfig:"WEBHOOK_URL"`
	Secret     string
}

type Service struct {
	config Config

	cachePrefix string
	cacheTTL    time.Duration
}

func (s *Service) InventoryAsset(ctx context.Context, steamID string) ([]steam.Asset, error) {
	// pull raw data from local cache
	return nil, nil
}

func (s *Service) SaveInventory(ctx context.Context, steamID string, r io.ReadCloser) error {
	var inventory Inventory
	if err := fastjson.NewDecoder(r).Decode(&inventory); err != nil {
		return fmt.Errorf("could not parse json form: %s", err)
	}
	defer func() {
		if err := r.Close(); err != nil {
			fmt.Printf("phantasm save inventory close reader: %s", err)
		}
	}()

	fmt.Printf("saving inventory %s\n", steamID)
	fmt.Printf("caching prefix %s\n", s.cachePrefix)
	fmt.Printf("caching ttl %s\n", s.cacheTTL)

	k := s.cacheKey(steamID)
	if err := localcache.Set(k, inventory, s.cacheTTL); err != nil {
		return fmt.Errorf("set cache: %s %s", k, err)
	}
	return nil
}

func (s *Service) cacheKey(steamID string) string {
	return fmt.Sprintf("%s_%s", s.cachePrefix, steamID)
}

func NewService(config Config) *Service {
	return &Service{
		config:      config,
		cachePrefix: "phantasm",
		cacheTTL:    time.Hour,
	}
}
