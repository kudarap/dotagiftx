package worker

import (
	"context"
	"fmt"
	"time"
)

type Trainee struct {
	ctr int
}

func (t *Trainee) ID() string {
	return "trainee"
}

func (t *Trainee) Run(ctx context.Context) error {
	fmt.Println("WORKER IS NOW TRAINING! ", t.ctr)
	t.ctr++
	return nil
}

func (t *Trainee) Interval() time.Duration {
	return time.Minute
}
