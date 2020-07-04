package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// Client represents Redis database client.
type Client struct {
	db *redis.Client
}

// New returns a new Redis client.
func New(addr, pass string, db int) (*Client, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	if err := c.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Client{c}, nil
}

// Close closes database client connection.
func (c *Client) Close() error {
	return c.db.Close()
}

func (c *Client) Set(key string, val interface{}, expr time.Duration) error {
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
