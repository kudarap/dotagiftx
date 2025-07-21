package rethink

import (
	"context"
	"log"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/xerrors"
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

func (s *taskStorage) Get(ctx context.Context) (*dotagiftx.Task, error) {
	res, err := s.List(ctx, 1)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}

func (s *taskStorage) List(ctx context.Context, limit int) ([]dotagiftx.Task, error) {
	q := s.table().GetAllByIndex(tableTaskFieldStatus, dotagiftx.TaskStatusPending).
		OrderBy(tableTaskFieldPriority).Limit(limit)

	var res []dotagiftx.Task
	if err := s.db.list(q, &res); err != nil {
		return nil, xerrors.New(dotagiftx.StorageUncaughtErr, err)
	}
	return res, nil
}

func (s *taskStorage) Update(ctx context.Context, in dotagiftx.Task) error {
	in.Retry++
	err := s.db.update(s.table().Get(in.ID).Update(in))
	if err != nil {
		return xerrors.New(dotagiftx.StorageUncaughtErr, err)
	}

	return nil
}

func (s *taskStorage) Queue(ctx context.Context, p dotagiftx.TaskPriority, t dotagiftx.TaskType, payload interface{}) (id string, err error) {
	n := now()
	id, err = s.db.insert(s.table().Insert(dotagiftx.Task{
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

	if err := c.autoIndex(tableTask, dotagiftx.Task{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableTask, err)
	}

	return &taskStorage{c}
}

func (s *taskStorage) table() r.Term {
	return r.Table(tableTask)
}
