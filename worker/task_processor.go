package worker

import (
	"context"
	"time"

	"github.com/kudarap/dotagiftx/core"
)

type TaskProcessor struct {
	queue taskQueue
	rate  time.Duration
}

func (p *TaskProcessor) Run() error {
	return nil
}

type taskQueue interface {
	Get(ctx context.Context) (core.Task, error)
	Update(ctx context.Context, status core.TaskStatus) error
}
