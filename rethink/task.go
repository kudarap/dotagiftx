package rethink

import (
	"context"
	"log"

	dotagiftx2 "github.com/kudarap/dotagiftx/dotagiftx"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	tableTask              = "task"
	tableTaskFieldStatus   = "status"
	tableTaskFieldPriority = "priority"
)

type TaskRepository struct {
	db *Client
}

func (s *TaskRepository) Get(ctx context.Context) (*dotagiftx2.Task, error) {
	res, err := s.List(ctx, 1)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}

func (s *TaskRepository) List(ctx context.Context, limit int) ([]dotagiftx2.Task, error) {
	q := s.table().GetAllByIndex(tableTaskFieldStatus, dotagiftx2.TaskStatusPending).
		OrderBy(tableTaskFieldPriority).Limit(limit)

	var res []dotagiftx2.Task
	if err := s.db.list(ctx, q, &res); err != nil {
		return nil, dotagiftx2.NewXError(dotagiftx2.StorageUncaughtErr, err)
	}
	return res, nil
}

func (s *TaskRepository) Update(ctx context.Context, in dotagiftx2.Task) error {
	in.Retry++
	err := s.db.update(ctx, s.table().Get(in.ID).Update(in))
	if err != nil {
		return dotagiftx2.NewXError(dotagiftx2.StorageUncaughtErr, err)
	}

	return nil
}

func (s *TaskRepository) Queue(ctx context.Context, p dotagiftx2.TaskPriority, t dotagiftx2.TaskType, payload any) (id string, err error) {
	n := now()
	id, err = s.db.insert(ctx, s.table().Insert(dotagiftx2.Task{
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

func NewQueue(c *Client) *TaskRepository {
	ctx := context.Background()
	if err := c.autoMigrate(ctx, tableTask); err != nil {
		log.Fatalf("could not create %s table: %s", tableTask, err)
	}

	if err := c.autoIndex(ctx, tableTask, dotagiftx2.Task{}); err != nil {
		log.Fatalf("could not create index on %s table: %s", tableTask, err)
	}

	return &TaskRepository{c}
}

func (s *TaskRepository) table() r.Term {
	return r.Table(tableTask)
}
