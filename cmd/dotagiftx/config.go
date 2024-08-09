package main

import (
	"github.com/kudarap/dotagiftx/gokit/log"
	"github.com/kudarap/dotagiftx/paypal"
	"github.com/kudarap/dotagiftx/redis"
	"github.com/kudarap/dotagiftx/rethink"
	"github.com/kudarap/dotagiftx/steam"
)

type Config struct {
	SigKey      string
	DivineKey   string
	Prod        bool
	Addr        string
	AppHost     string
	ApiHost     string
	SpanEnabled bool `envconfig:"SPAN_ENABLED"`
	Upload      struct {
		Path  string
		Size  int
		Types []string
	}
	Rethink           rethink.Config
	Redis             redis.Config
	Steam             steam.Config
	Paypal            paypal.Config
	Log               log.Config
	DiscordWebhookURL string `envconfig:"DISCORD_WEBHOOK_URL"`
}
