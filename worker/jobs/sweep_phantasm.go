package jobs

import (
	"context"
	"time"

	"github.com/kudarap/dotagiftx/logging"
)

type CleanPhantasmCache struct {
	service     phantasmCacheCleaner
	cacheMaxAge time.Duration

	name     string
	interval time.Duration
	logger   logging.Logger
}

func NewSweepPhantasmCache(cleaner phantasmCacheCleaner, lg logging.Logger) *CleanPhantasmCache {
	return &CleanPhantasmCache{
		service:     cleaner,
		cacheMaxAge: time.Hour * 24 * 30, // 30 days
		name:        "clean_phantasm_cache",
		interval:    time.Hour * 24,
		logger:      lg,
	}
}

func (c *CleanPhantasmCache) String() string { return c.name }

func (c *CleanPhantasmCache) Interval() time.Duration { return c.interval }

func (c *CleanPhantasmCache) Run(ctx context.Context) error {
	c.logger.Println("cleaning phantasm cache older than", c.cacheMaxAge)
	if err := c.service.CleanLocalCache(ctx, c.cacheMaxAge); err != nil {
		return err
	}
	c.logger.Println("phantasm cache cleaned!")
	return nil
}

type phantasmCacheCleaner interface {
	CleanLocalCache(ctx context.Context, maxAge time.Duration) error
}
