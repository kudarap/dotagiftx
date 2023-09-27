package rethink

import (
	"log"

	"github.com/kudarap/dotagiftx/core"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	tableTask = "task"
)

type taskStorage struct {
	db *Client
}

func NewQueue(c *Client) *taskStorage {
	if err := c.autoMigrate(tableTask); err != nil {
		log.Fatalf("could not create %s table: %s", tableTrack, err)
	}

	if err := c.autoIndex(tableTask, core.Track{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableMarket, err)
	}

	return &taskStorage{c}
}

func (s *taskStorage) VerifyDelivery(marketID string) {
	var t core.Task
	t.Type = core.TaskTypeVerifyDelivery
	t.Payload = marketID
	t.CreatedAt = now()
	_, err := s.db.insert(s.table().Insert(t))
	if err != nil {
		log.Println("ERR TASK VerifyDelivery", err)
	}
}

func (s *taskStorage) VerifyInventory(userID string) {
	var t core.Task
	t.Type = core.TaskTypeVerifyInventory
	t.Payload = userID
	t.CreatedAt = now()
	_, err := s.db.insert(s.table().Insert(t))
	if err != nil {
		log.Println("ERR TASK VerifyInventory", err)
	}
}

func (s *taskStorage) table() r.Term {
	return r.Table(tableTask)
}
