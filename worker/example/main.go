package main

import (
	"context"
	"log"
	"time"

	"github.com/kudarap/dotagiftx/worker"
)

func main() {
	w := worker.New(
		NewTraineeJob("grind beans"),
		NewTraineeJob("boil some good fucking watter"),
		NewTraineeJob("get a cup"),
	)
	go w.Start()

	//w.AddJob(NewTraineeJob("make it fly"))
	time.Sleep(time.Second * 2)
	//w.AddJob(NewTraineeRunOnceJob("drink once"))

	// Initiates early termination will finish the remaining jobs
	//time.Sleep(time.Minute)
	if err := w.Stop(); err != nil {
		log.Println("could not stop worker:", err)
	}
}

// Trainee represents a sample job that implements worker.Job
type Trainee struct {
	ctr      int
	name     string
	interval time.Duration
}

func NewTraineeJob(name string) *Trainee {
	return &Trainee{0, name, time.Second * 5}
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
