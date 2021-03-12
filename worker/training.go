package worker

import (
	"context"
	"fmt"
	"time"
)

type Trainee struct {
	ctr int
}

func NewTraineeJob() *Trainee {
	return &Trainee{}
}

func (t *Trainee) ID() string {
	return "trainee"
}

func (t *Trainee) Run(ctx context.Context) error {
	fmt.Println("WORKER IS NOW TRAINING! started", t.ctr)
	time.Sleep(time.Second * 10)
	fmt.Println("WORKER IS NOW TRAINING! ended", t.ctr)
	t.ctr++
	return nil
}

func (t *Trainee) Interval() time.Duration {
	return time.Second / 5
}
