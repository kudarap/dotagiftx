package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/kudarap/dotagiftx/clickhouse"
	"github.com/kudarap/dotagiftx/file"
	"github.com/kudarap/dotagiftx/github"
	"github.com/kudarap/dotagiftx/logging"
	"github.com/kudarap/dotagiftx/paypal"
	"github.com/kudarap/dotagiftx/phantasm"
	"github.com/kudarap/dotagiftx/redis"
	"github.com/kudarap/dotagiftx/rethink"
	"github.com/kudarap/dotagiftx/steam"
)

// EnvPrefix default env prefix APP.
var EnvPrefix = "APP"

type Config struct {
	SigKey              string
	DivineKey           string
	Prod                bool
	Addr                string
	AppHost             string `envconfig:"APP_HOST"`
	ApiHost             string `envconfig:"API_HOST"`
	SpanEnabled         bool   `envconfig:"SPAN_ENABLED"`
	StatsCaptureEnabled bool   `envconfig:"STATS_CAPTURE_ENABLED"`
	Upload              file.Config
	AllowedImageSources []string `envconfig:"ALLOWED_IMAGE_SOURCES"`
	Rethink             rethink.Config
	Redis               redis.Config
	ClickHouse          clickhouse.Config
	Steam               steam.Config
	Paypal              paypal.Config
	Log                 logging.Config
	Phantasm            phantasm.Config
	DiscordWebhookURL   string `envconfig:"DISCORD_WEBHOOK_URL"`
	Github              github.Config
	AuthSessionTTL      time.Duration
}

// Load parses .env values into a struct.
func Load(conf *Config) error {
	// Load env file.
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("could not load config: %s", err)
	}
	// Bind env values.
	if err := envconfig.Process(EnvPrefix, conf); err != nil {
		return fmt.Errorf("could not process config: %s", err)
	}
	return nil
}
