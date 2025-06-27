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
	phantasmRetryAfterKey = phantasmKey + "/retryafter"
	phantasmCooldownKey   = phantasmKey + "/cooldown"
	phantasmInvHashKey    = phantasmKey + "/invhash"
)

func (c *Client) RetryCooldown(ctx context.Context, crawlID, steamID string) (bool, error) {
	key := fmt.Sprintf("%s/%s/%s", phantasmRetryAfterKey, crawlID, steamID)

	v, err := c.db.Get(ctx, key).Bool()
	if err != nil && !errors.Is(err, redis.Nil) {
		return v, err
	}
	return v, nil
}

func (c *Client) SetRetryCooldown(ctx context.Context, crawlID, steamID string, ttl time.Duration) error {
	key := fmt.Sprintf("%s/%s/%s", phantasmRetryAfterKey, crawlID, steamID)
	return c.db.Set(ctx, key, true, ttl).Err()
}

func (c *Client) CrawlerCooldown(ctx context.Context, crawlID string) (bool, error) {
	key := fmt.Sprintf("%s/%s", phantasmCooldownKey, crawlID)

	v, err := c.db.Get(ctx, key).Bool()
	if err != nil && !errors.Is(err, redis.Nil) {
		return v, err
	}
	return v, nil
}

func (c *Client) SetCrawlerCooldown(ctx context.Context, crawlID string, ttl time.Duration) error {
	key := fmt.Sprintf("%s/%s", phantasmCooldownKey, crawlID)
	return c.db.Set(ctx, key, true, ttl).Err()
}

func (c *Client) InventoryHash(ctx context.Context, steamID string) (hash string, error error) {
	key := fmt.Sprintf("%s/%s", phantasmInvHashKey, steamID)
	v, err := c.db.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return v, err
	}
	return v, nil
}

func (c *Client) SetInventoryHash(ctx context.Context, steamID, hash string, ttl time.Duration) error {
	key := fmt.Sprintf("%s/%s", phantasmInvHashKey, steamID)
	return c.db.Set(ctx, key, hash, ttl).Err()
}
