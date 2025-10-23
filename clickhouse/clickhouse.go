package clickhouse

import (
	"context"
	"errors"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type Client struct {
	conn driver.Conn
	conf Config
}

func New(conf Config) (*Client, error) {
	ctx := context.Background()
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"localhost:9000"},
		Auth: clickhouse.Auth{},
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

	return &Client{conn, conf}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

type Config struct {
	Addr string
}
