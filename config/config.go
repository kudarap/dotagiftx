package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/kudarap/dotagiftx/file"
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
	SigKey            string
	DivineKey         string
	Prod              bool
	Addr              string
	AppHost           string
	ApiHost           string
	SpanEnabled       bool `envconfig:"SPAN_ENABLED"`
	Upload            file.Config
	Rethink           rethink.Config
	Redis             redis.Config
	Steam             steam.Config
	Paypal            paypal.Config
	Log               logging.Config
	Phantasm          phantasm.Config
	DiscordWebhookURL string `envconfig:"DISCORD_WEBHOOK_URL"`
}

// Load parses .env values into a struct.
func Load(conf *Config) error {
	// Load env file.
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("could not load config: %s", err)
	}
	// Bind env values.
	if err := envconfig.Process(EnvPrefix, conf); err != nil {
		return fmt.Errorf("could not process config: %s", err)
	}
	return nil
}
