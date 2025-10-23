package clickhouse

import (
	"context"
	"errors"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type Client struct {
	db   driver.Conn
	conf Config
}

func New(conf Config) (*Client, error) {
	ctx := context.Background()
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{conf.Addr},
		Auth: clickhouse.Auth{
			Database: conf.Db,
			Username: conf.User,
			Password: conf.Pass,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("connect: %s", err)
	}

	if err = conn.Ping(ctx); err != nil {
		var exception *clickhouse.Exception
		if errors.As(err, &exception) {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, fmt.Errorf("ping: %s", err)
	}

	if err = autoMigration(ctx, conn); err != nil {
		return nil, fmt.Errorf("auto migration: %s", err)
	}

	return &Client{conn, conf}, nil
}

func (c *Client) Close() error {
	return c.db.Close()
}

type Config struct {
	Addr string
	Db   string
	User string
	Pass string
}
