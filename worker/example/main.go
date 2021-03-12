package main

import (
	"log"
	"time"

	"github.com/kudarap/dotagiftx/worker"
)

func main() {

	w := worker.New(
		worker.NewTraineeJob("KARLINGKOMORO"),
		//worker.NewTraineeJob("KUDARAP"),
		//worker.NewTraineeJob("MOMO"),
	)
	go w.Start()

	time.Sleep(time.Second * 9)
	w.AddJob(worker.NewTraineeJob("IM ON THE FLY JOB"))
	w.AddJob(worker.NewTraineeRunOnceJob("IM A RUN ONCE JOB"))

	// Initiates early termination will finish the remaining jobs
	time.Sleep(time.Minute)
	if err := w.Stop(); err != nil {
		log.Println("could not stop worker:", err)
	}
}
