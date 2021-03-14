package worker

import (
	"context"
	"sync"
	"time"

	"github.com/kudarap/dotagiftx/gokit/log"
)

// Worker represents worker handling and running tasks.
type Worker struct {
	wg     sync.WaitGroup
	quit   chan struct{}
	queue  chan Job
	jobs   []Job
	closed bool

	logger log.Logger
}

// New create new instance of a worker with a given jobs.
func New(jobs ...Job) *Worker {
	w := &Worker{}
	w.queue = make(chan Job, len(jobs))
	w.quit = make(chan struct{})
	w.jobs = jobs

	w.logger = log.Default()
	return w
}

// SetLogger overrides default logger.
func (w *Worker) SetLogger(l log.Logger) {
	w.logger = l
}

// Start initiates worker to start running the jobs.
//
// All assigned jobs will be run concurrently.
func (w *Worker) Start() {
	w.logger.Infof("running %d jobs gracefully...", len(w.jobs))

	ctx := context.Background()

	// Run registered jobs gracefully on start to
	// prevents initial burst of server load.
	for _, jj := range w.jobs {
		//var delay time.Duration
		//go func(j Job) {
		//w.runner(ctx, j, false)
		//delay += j.Interval()
		//time.Sleep(delay)
		//}(jj)

		//w.runner(ctx, jj, true)

		go w.queueJob(jj)
	}

	// Handles job queueing and termination.
	for {
		select {
		// Job queue is now closed and will not run jobs anymore.
		// Queued jobs will be terminated.
		case <-w.quit:
			return
		case job, ok := <-w.queue:
			if !ok {
				return
			}

			w.runner(ctx, job, false)
			// Enable this of you want multi-tasking worker.
			//go w.runner(ctx, job, false)
		}
	}
}

// AddJob registers a new job while the worker is running.
func (w *Worker) AddJob(j Job) {
	//w.jobs = append(w.jobs, j)
	w.queueJob(j)
}

// runner process the job and will re-queue them when recurring job.
func (w *Worker) runner(ctx context.Context, job Job, once bool) {
	w.logger.Infof("RUNN job:%s", job)
	w.wg.Add(1)

	if err := job.Run(ctx); err != nil {
		w.logger.Errorf("ERRO job:%s - %s", job, err)
	}
	w.logger.Infof("DONE job:%s", job)
	w.wg.Done()

	// Worker queue is now closed.
	if w.closed {
		w.logger.Warnf("SKIP job:%s queue is closed", job)
		return
	}

	// Job that has non-zero interval value means its a recurring job
	// and will be re-queued after its rest duration.
	if once {
		return
	}

	// Determines if the job is run-once by interval value is zero,
	rest := job.Interval()
	if rest == 0 {
		return
	}

	w.logger.Infof("REST job:%s will re-queue in %s", job, rest)
	time.Sleep(rest)
	w.queueJob(job)
}

func (w *Worker) queueJob(j Job) {
	if w.closed {
		w.logger.Warnf("SKIP job:%s queue is closed", j)
		return
	}

	w.logger.Printf("TODO job:%s", j)
	w.queue <- j
}

// Stop will stop accepting job and wait for processing job to finish.
func (w *Worker) Stop() error {
	w.logger.Infof("stopping and waiting for jobs to finish...")
	w.quit <- struct{}{}
	close(w.queue)
	w.closed = true
	w.wg.Wait()
	w.logger.Infof("all jobs done!")
	return nil
}

// RunOnce will queue the job but it will not register to worker's jobs,
// since its one-time job it will ignore the Interval() value.
//
// If you want a recurring Job, you must register it worker constructor.
// This method is design for on-the-fly indexing or on demand delivery
// verification.
//
// DEPRECATED
func (w *Worker) RunOnce(j Job) {
	w.queueJob(j)
}
