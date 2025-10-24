package rethink

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func (c *Client) ListenChangeFeed(table string, exec func(next, prev []byte) error) error {
	feed, err := newChangeFeed(c.db, table, exec)
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

func newChangeFeed(db *r.Session, table string, exec func(next, prev []byte) error) (*changeFeed, error) {
	t := r.Table(table).Changes()
	cursor, err := t.Run(db)
	if err != nil {
		return nil, err
	}

	var feed changeFeed
	feed.ch = make(chan map[string]any, 10000)
	feed.closer = make(chan bool)
	feed.cursor = cursor

	logrus.Info(table, "change feed started")
	go func() {
		feed.cursor.Listen(feed.ch)
		for {
			select {
			case <-feed.closer:
				logrus.Info(table, "change feed closed")
				return

			case event := <-feed.ch:
				next, err := json.Marshal(event["new_val"])
				if err != nil {
					logrus.Errorf("could not marshal new_val: %s", err)
					continue
				}
				prev, err := json.Marshal(event["old_val"])
				if err != nil {
					logrus.Errorf("could not marshal new_val: %s", err)
					continue
				}
				if err = exec(next, prev); err != nil {
					logrus.Errorf("could not process change: %s", err)
				}
			}
		}
	}()

	return &feed, nil
}
