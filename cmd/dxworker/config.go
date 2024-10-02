package main

import (
	"github.com/kudarap/dotagiftx/gokit/log"
	"github.com/kudarap/dotagiftx/redis"
	"github.com/kudarap/dotagiftx/rethink"
	"github.com/kudarap/dotagiftx/steam"
)

type Config struct {
	SigKey      string
	Prod        bool
	Addr        string
	SpanEnabled bool `envconfig:"SPAN_ENABLED"`
	Rethink     rethink.Config
	Redis       redis.Config
	Steam       steam.Config
	Log         log.Config
}
