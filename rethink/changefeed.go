package rethink

import (
	"context"
	"encoding/json"
	"log/slog"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func (c *Client) ListenChangeFeed(ctx context.Context, table string, exec func(prev, next []byte) error) error {
	feed, err := newChangeFeed(ctx, c.db, table, exec)
	if err != nil {
		return err
	}

	c.changeFeeds = append(c.changeFeeds, feed)
	return nil
}

type changeFeed struct {
	ch     chan map[string]any
	closer chan bool
	cursor *r.Cursor
}

func (f *changeFeed) close() error {
	f.closer <- true
	return f.cursor.Close()
}

func newChangeFeed(ctx context.Context, db *r.Session, table string, exec func(prev, next []byte) error) (*changeFeed, error) {
	t := r.Table(table).Changes()
	cursor, err := t.Run(db, r.RunOpts{Context: ctx})
	if err != nil {
		return nil, err
	}

	var feed changeFeed
	feed.ch = make(chan map[string]any, 10000)
	feed.closer = make(chan bool)
	feed.cursor = cursor

	slog.InfoContext(ctx, "change feed started", "table", table)
	go func() {
		feed.cursor.Listen(feed.ch)
		for {
			select {
			case <-feed.closer:
				slog.InfoContext(ctx, "change feed closed", "table", table)
				return

			case event := <-feed.ch:
				var oldVal, newVal []byte
				if raw := event["old_val"]; raw != nil {
					oldVal, err = json.Marshal(raw)
					if err != nil {
						slog.ErrorContext(ctx, "could not marshal old_val", "error", err)
						continue
					}
				}
				if raw := event["new_val"]; raw != nil {
					newVal, err = json.Marshal(raw)
					if err != nil {
						slog.ErrorContext(ctx, "could not marshal new_val", "error", err)
						continue
					}
				}
				if err = exec(oldVal, newVal); err != nil {
					slog.ErrorContext(ctx, "could not process change", "error", err)
				}
			}
		}
	}()

	return &feed, nil
}
