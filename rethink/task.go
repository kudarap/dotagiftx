package rethink

import (
	"context"
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

func (s *taskStorage) Queue(ctx context.Context, p core.TaskPriority, t core.TaskType, payload interface{}) (id string, err error) {
	n := now()
	id, err = s.db.insert(s.table().Insert(core.Task{
		Status:    0,
		Priority:  p,
		Type:      t,
		Payload:   payload,
		CreatedAt: n,
		UpdatedAt: n,
	}))
	if err != nil {
		return "", err
	}
	return id, nil
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

func (s *taskStorage) table() r.Term {
	return r.Table(tableTask)
}
