package worker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// JobID represents identification for a Job.
type JobID string

// Job provides process and information for the job.
type Job interface {
	ID() string
	Task
}

type Task interface {
	Run(context.Context) error
	Interval() time.Duration
}

// Worker represents worker handling and running tasks.
type Worker struct {
	wg    sync.WaitGroup
	quit  chan struct{}
	queue chan Job
	jobs  []Job
}

// New create new instance of a worker with a given jobs.
func New(jobs ...Job) *Worker {
	w := &Worker{}
	w.queue = make(chan Job, 1)
	w.quit = make(chan struct{})
	w.jobs = jobs

	return w
}

// Start initiates worker to start running the jobs.
//
// All assigned jobs will be run concurrently.
func (w *Worker) Start() error {
	w.logger("running", len(w.jobs), "jobs...")

	go func() {
		for _, jj := range w.jobs {
			w.queueJob(jj)
		}
	}()

	ctx := context.Background()
	for {
		select {
		case <-w.quit:
			return nil
		case job := <-w.queue:
			go w.runner(ctx, job)
		}
	}
}

func (w *Worker) RunOnce() {

}

// runner process the job and will re-queue them when recurring job.
func (w *Worker) runner(ctx context.Context, task Job) {
	w.wg.Add(1)

	w.logger(fmt.Sprintf("job:%s recv", task))
	if err := task.Run(ctx); err != nil {
		w.logger(fmt.Sprintf("job:%s error! %s", task, err))
	}
	w.logger(fmt.Sprintf("job:%s done!", task))
	w.wg.Done()

	// Determines if the job is run-once by interval value is zero.
	rest := task.Interval()
	if rest == 0 {
		return
	}
	// Job that has non-zero interval value means its a recurring job
	// and will be re-queued after its rest duration.
	time.Sleep(rest)
	w.queueJob(task)
}

func (w *Worker) queueJob(j Job) {
	w.queue <- j
	w.logger(fmt.Sprintf("job:%s queued", j))
}

// Stop will stop accepting job and wait for processing job to finish.
func (w *Worker) Stop() error {
	w.logger("stopping and waiting for jobs finish...")
	w.quit <- struct{}{}
	w.wg.Wait()
	w.logger("all jobs done!")
	return nil
}

func (w *Worker) logger(v ...interface{}) {
	v = append([]interface{}{"[worker]"}, v...)
	log.Println(v...)
}
