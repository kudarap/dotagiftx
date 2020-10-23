package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigFastest

// Config represents redis database config.
type Config struct {
	Addr string
	Db   int
	Pass string
}

var ctx = context.Background()

// Client represents Redis database client.
type Client struct {
	db *redis.Client
}

// New returns a new Redis client.
func New(c Config) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Pass,
		DB:       c.Db,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Client{rdb}, nil
}

// Close closes database client connection.
func (c *Client) Close() error {
	return c.db.Close()
}

func (c *Client) Set(key string, val interface{}, expr time.Duration) error {
	// Skip caching when key and value is empty.
	if key == "" || val == nil {
		return nil
	}

	b, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return c.db.Set(ctx, key, string(b), expr).Err()
}

func (c *Client) Get(key string) (val string, err error) {
	res, err := c.db.Get(ctx, key).Result()
	if err == redis.Nil {
		err = nil
		return
	} else if err != nil {
		return
	}

	return res, nil
}
