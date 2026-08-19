package jobs

import (
	"context"
	"log/slog"
	"time"
)

type CleanPhantasmCache struct {
	service     phantasmCacheCleaner
	cacheMaxAge time.Duration

	name     string
	interval time.Duration
	logger   *slog.Logger
}

func NewSweepPhantasmCache(cleaner phantasmCacheCleaner, lg *slog.Logger) *CleanPhantasmCache {
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
	c.logger.Info("cleaning phantasm cache older than", "max_age", c.cacheMaxAge)
	if err := c.service.CleanLocalCache(ctx, c.cacheMaxAge); err != nil {
		return err
	}
	c.logger.Info("phantasm cache cleaned!")
	return nil
}

type phantasmCacheCleaner interface {
	CleanLocalCache(ctx context.Context, maxAge time.Duration) error
}
