package main

import (
	"time"

	"github.com/kudarap/dotagiftx/worker"
)

func main() {

	w := worker.New(worker.NewTraineeJob())
	go w.Start()
	time.Sleep(time.Second / 5)
	w.Stop()
}
