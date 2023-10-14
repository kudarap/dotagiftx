package rethink

import (
	"context"
	"log"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/errors"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	tableTask              = "task"
	tableTaskFieldStatus   = "status"
	tableTaskFieldPriority = "priority"
)

type taskStorage struct {
	db *Client
}

func (s *taskStorage) Get(ctx context.Context) (*core.Task, error) {
	res, err := s.List(ctx, 1)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}

func (s *taskStorage) List(ctx context.Context, limit int) ([]core.Task, error) {
	q := s.table().GetAllByIndex(tableTaskFieldStatus, core.TaskStatusPending).
		OrderBy(tableTaskFieldPriority).Limit(limit)

	var res []core.Task
	if err := s.db.list(q, &res); err != nil {
		return nil, errors.New(core.StorageUncaughtErr, err)
	}
	return res, nil
}

func (s *taskStorage) Update(ctx context.Context, in core.Task) error {
	in.Retry++
	err := s.db.update(s.table().Get(in.ID).Update(in))
	if err != nil {
		return errors.New(core.StorageUncaughtErr, err)
	}

	return nil
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
		log.Fatalf("could not create %s table: %s", tableTask, err)
	}

	if err := c.autoIndex(tableTask, core.Task{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableTask, err)
	}

	return &taskStorage{c}
}

func (s *taskStorage) table() r.Term {
	return r.Table(tableTask)
}
