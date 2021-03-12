package worker

import (
	"context"
	"time"
)

type Trainee struct {
	ctr      int
	name     string
	interval time.Duration
}

func NewTraineeJob(name string) *Trainee {
	return &Trainee{0, name, time.Second / 5}
}

func NewTraineeRunOnceJob(name string) *Trainee {
	return &Trainee{0, name, 0}
}

func (t *Trainee) String() string { return t.name }

func (t *Trainee) Interval() time.Duration { return t.interval }

func (t *Trainee) Run(ctx context.Context) error {
	t.ctr++
	//log.Println(t.name, "started working on #", t.ctr)
	time.Sleep(time.Second * time.Duration(len(t.name)))
	//log.Println(t.name, "finished working on #", t.ctr)
	return nil
}
