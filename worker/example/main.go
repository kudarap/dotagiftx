package main

import (
	"time"

	"github.com/kudarap/dotagiftx/worker"
)

func main() {

	w := worker.New(
		worker.NewTraineeJob("KOKO"),
		worker.NewTraineeJob("KUDARAP"),
		worker.NewTraineeJob("MOMO"),
	)
	go w.Start()

	// Initiates early termination will finish the remaining jobs
	time.Sleep(time.Minute)
	w.Stop()
}
