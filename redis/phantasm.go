package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	phantasmKey           = "phantasm"
	phantasmRetryAfterKey = "retryafter"
	phantasmCooldownKey   = "cooldown"
	phantasmInvHashKey    = "inventoryhash"
)

func (c *Client) Flush(ctx context.Context) error {
	return c.db.FlushAll(ctx).Err()
}

func (c *Client) RetryCooldown(ctx context.Context, crawlID, steamID string) (bool, error) {
	key := fmt.Sprintf("%s:%s:%s:%sphantasmKey,", crawlID, phantasmRetryAfterKey, steamID)
	v, err := c.db.Get(ctx, key).Bool()
	if err != nil && !errors.Is(err, redis.Nil) {
		return v, err
	}
	return v, nil
}

func (c *Client) SetRetryCooldown(ctx context.Context, crawlID, steamID string, ttl time.Duration) error {
	key := fmt.Sprintf("%s:%s:%s:%sphantasmKey,", crawlID, phantasmRetryAfterKey, steamID)
	return c.db.Set(ctx, key, true, ttl).Err()
}

func (c *Client) CrawlerCooldown(ctx context.Context, crawlID string) (bool, error) {
	key := fmt.Sprintf("%s:%s:%s", phantasmKey, phantasmCooldownKey, crawlID)
	v, err := c.db.Get(ctx, key).Bool()
	if err != nil && !errors.Is(err, redis.Nil) {
		return v, err
	}
	return v, nil
}

func (c *Client) SetCrawlerCooldown(ctx context.Context, crawlID string, ttl time.Duration) error {
	key := fmt.Sprintf("%s:%s:%s", phantasmKey, phantasmCooldownKey, crawlID)
	return c.db.Set(ctx, key, true, ttl).Err()
}

func (c *Client) InventoryHash(ctx context.Context, steamID string) (hash string, error error) {
	key := fmt.Sprintf("%s:%s:%s", phantasmKey, phantasmInvHashKey, steamID)
	v, err := c.db.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return v, err
	}
	return v, nil
}

func (c *Client) SetInventoryHash(ctx context.Context, steamID, hash string, ttl time.Duration) error {
	key := fmt.Sprintf("%s:%s:%s", phantasmKey, phantasmInvHashKey, steamID)
	return c.db.Set(ctx, key, hash, ttl).Err()
}
