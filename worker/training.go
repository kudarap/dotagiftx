package worker

import (
	"context"
	"log"
	"time"
)

type Trainee struct {
	ctr  int
	name string
}

func NewTraineeJob(name string) *Trainee {
	return &Trainee{0, name}
}

func (t *Trainee) ID() string {
	return "trainee_" + t.name
}

func (t *Trainee) Run(ctx context.Context) error {
	t.ctr++
	log.Println(t.name, "started working on #", t.ctr)
	time.Sleep(time.Second * time.Duration(len(t.name)))
	log.Println(t.name, "finished working on #", t.ctr)
	return nil
}

func (t *Trainee) Interval() time.Duration {
	return time.Second / 5
}
