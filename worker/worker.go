package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kudarap/dotagiftx/gokit/log"
	"github.com/kudarap/dotagiftx/tracing"
)

// Worker represents worker handling and running tasks.
type Worker struct {
	wg       sync.WaitGroup
	quit     chan struct{}
	queue    chan Job
	jobs     []Job
	closed   bool
	taskProc *TaskProcessor

	logger log.Logger
	tracer *tracing.Tracer
}

// JobID represents identification for a Job.
type JobID string

// Job provides process and information for the job.
type Job interface {
	// String returns job reference or name.
	String() string

	// Interval returns sleep duration before re-queueing.
	//
	// Returning a zero value will consider run-once job
	// and will NOT be re-queued.
	Interval() time.Duration

	// Run process the task of the job.
	//
	// Recurring Job will not stop when an error occurred.
	Run(context.Context) error
}

// New create new instance of a worker with a given jobs.
func New(tp *TaskProcessor, jobs ...Job) *Worker {
	w := &Worker{}
	w.queue = make(chan Job, len(jobs))
	w.quit = make(chan struct{})
	w.jobs = jobs
	w.taskProc = tp
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
	w.logger.Infof("running task processor...")
	go w.taskProc.Run(&w.wg)

	w.logger.Infof("running jobs...")

	ctx := context.Background()

	// Queue initial registered jobs.
	for _, jj := range w.jobs {
		w.queueJob(jj, true)
	}

	// Handles job queueing and termination.
	for {
		select {
		// Waits the running job to finish and stops the worker and skip all queued jobs.
		case <-w.quit:
			return
		case job := <-w.queue:
			if w.closed {
				w.logger.Warnf("SKIP job:%s queue is closed", job)
				continue
			}

			// Runner will block until done making it a single tasking worker.
			w.runner(ctx, job)

			// Enable this of you want multitasking worker.
			//go w.runner(ctx, job, false)
		}
	}
}

// AddJob registers a new job while the worker is running.
func (w *Worker) AddJob(j Job) {
	//w.jobs = append(w.jobs, j)
	w.queueJob(j, true)
}

// runner process the job and will re-queue them when recurring job.
func (w *Worker) runner(ctx context.Context, job Job) {
	if w.tracer != nil {
		span := w.tracer.StartSpan(fmt.Sprintf("job-%s", job))
		defer func() {
			span.End()
		}()
	}

	w.logger.Infof("RUNN job:%s", job)
	w.wg.Add(1)

	if err := job.Run(ctx); err != nil {
		w.logger.Errorf("ERRO job:%s - %s", job, err)
	}
	w.logger.Infof("DONE job:%s", job)
	w.wg.Done()

	// Worker queue is now closed.
	if w.closed {
		return
	}

	// Determines if the job is run-once by interval value is zero,
	rest := job.Interval()
	if rest == 0 {
		return
	}

	// Job that has non-zero interval value means it's a recurring job
	// and will be re-queued after its rest duration
	w.logger.Infof("REST job:%s will re-queue in %s", job, rest)
	w.queueJob(job, false)
}

// queueJob handles job whether it should be queued immediately or
// standby aside and wait for its next iteration.
//
// waiting jobs will be vanished into abyss when the worker is done.
func (w *Worker) queueJob(j Job, now bool) {
	if w.closed {
		w.logger.Warnf("SKIP job:%s queue is closed", j)
		return
	}

	// This will make the job to wait in goroutine before
	// adding itself to queue.
	//
	// This allows us to run another job from queue without
	// waiting for the interval to finish and making it
	// follow its own interval.
	go func() {
		if !now {
			time.Sleep(j.Interval())
		}
		w.logger.Printf("TODO job:%s", j)
		w.queue <- j
	}()
}

// Stop will stop accepting job and wait for processing job to finish.
func (w *Worker) Stop() error {
	w.logger.Infof("stopping and waiting for jobs to finish...")
	w.closed = true
	w.quit <- struct{}{}
	w.wg.Wait()
	w.logger.Infof("all jobs done!")
	return nil
}

// RunOnce will queue the job, but it will not register to worker's jobs,
// since its one-time job it will ignore the Interval() value.
//
// If you want a recurring Job, you must register it worker constructor.
// This method is design for on-the-fly indexing or on demand delivery
// verification.
//
// DEPRECATED
func (w *Worker) RunOnce(j Job) {
	w.queueJob(j, true)
}

func (w *Worker) SetTracer(t *tracing.Tracer) {
	w.tracer = t
}
