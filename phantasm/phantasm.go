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
	"time"

	"github.com/kudarap/dotagiftx/steam"
)

const (
	cacheExpr   = time.Hour * 24
	cachePrefix = "phantasm"
)

type Config struct {
	Addrs      []string
	WebhookURL string `envconfig:"WEBHOOK_URL"`
	Secret     string
}

type Service struct {
	config Config
}

func (s *Service) InventoryAsset(ctx context.Context, steamID string) ([]steam.Asset, error) {
	// pull raw data from local cache
	//
	return nil, nil
}
