package main

import (
	"github.com/kudarap/dota2giftables/gokit/logger"
	"github.com/kudarap/dota2giftables/rethink"
	"github.com/kudarap/dota2giftables/steam"
)

type (
	Config struct {
		SigKey  string
		Prod    bool
		Addr    string
		AppHost string
		ApiHost string

		Upload struct {
			Path  string
			Size  int
			Types []string
		}

		Rethink rethink.Config

		Steam steam.Config

		Log logger.Config
	}
)
