package rethink

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func (c *Client) ListenChangeFeed(table string, exec func([]byte) error) error {
	feed, err := newChangeFeed(c.db, table, exec)
	if err != nil {
		return err
	}

	c.changeFeeds = append(c.changeFeeds, feed)
	return nil
}

type changeFeed struct {
	cursor *r.Cursor
	closed chan bool
}

func (f *changeFeed) close() error {
	f.closed <- true
	return f.cursor.Close()
}

func newChangeFeed(db *r.Session, table string, fn func([]byte) error) (*changeFeed, error) {
	t := r.Table(table).Changes()
	cursor, err := t.Run(db)
	if err != nil {
		return nil, err
	}

	closed := make(chan bool, 1)
	feed := make(chan map[string]any)
	cursor.Listen(feed)
	logrus.Info(table, "change feed started")
	go func() {
		for {
			select {
			case <-closed:
				logrus.Info(table, "change feed closed")
				return

			case event := <-feed:
				b, err := json.Marshal(event["new_val"])
				if err != nil {
					logrus.Errorf("could not marshal new_val: %s", err)
				}
				if err = fn(b); err != nil {
					logrus.Errorf("could not process change: %s", err)
				}
			}
		}
	}()

	return &changeFeed{cursor, closed}, nil
}
