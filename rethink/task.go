package rethink

import (
	"context"
	"log"

	"github.com/kudarap/dotagiftx"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	tableTask              = "task"
	tableTaskFieldStatus   = "status"
	tableTaskFieldPriority = "priority"
)

type TaskStorage struct {
	db *Client
}

func (s *TaskStorage) Get(ctx context.Context) (*dotagiftx.Task, error) {
	res, err := s.List(ctx, 1)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}

func (s *TaskStorage) List(ctx context.Context, limit int) ([]dotagiftx.Task, error) {
	q := s.table().GetAllByIndex(tableTaskFieldStatus, dotagiftx.TaskStatusPending).
		OrderBy(tableTaskFieldPriority).Limit(limit)

	var res []dotagiftx.Task
	if err := s.db.list(ctx, q, &res); err != nil {
		return nil, dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}
	return res, nil
}

func (s *TaskStorage) Update(ctx context.Context, in dotagiftx.Task) error {
	in.Retry++
	err := s.db.update(ctx, s.table().Get(in.ID).Update(in))
	if err != nil {
		return dotagiftx.NewXError(dotagiftx.StorageUncaughtErr, err)
	}

	return nil
}

func (s *TaskStorage) Queue(ctx context.Context, p dotagiftx.TaskPriority, t dotagiftx.TaskType, payload any) (id string, err error) {
	n := now()
	id, err = s.db.insert(ctx, s.table().Insert(dotagiftx.Task{
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

func NewQueue(c *Client) *TaskStorage {
	ctx := context.Background()
	if err := c.autoMigrate(ctx, tableTask); err != nil {
		log.Fatalf("could not create %s table: %s", tableTask, err)
	}

	if err := c.autoIndex(ctx, tableTask, dotagiftx.Task{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableTask, err)
	}

	return &TaskStorage{c}
}

func (s *TaskStorage) table() r.Term {
	return r.Table(tableTask)
}
