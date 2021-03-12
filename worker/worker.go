package worker

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/kudarap/dotagiftx/gokit/logger"
	"github.com/sirupsen/logrus"
)

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

// Worker represents worker handling and running tasks.
type Worker struct {
	wg    sync.WaitGroup
	quit  chan struct{}
	queue chan Job
	jobs  []Job

	logger *logrus.Entry
}

// New create new instance of a worker with a given jobs.
func New(jobs ...Job) *Worker {
	w := &Worker{}
	w.queue = make(chan Job, 1)
	w.quit = make(chan struct{})
	w.jobs = jobs

	w.logger = logger.WithPrefix("worker")
	return w
}

func (w *Worker) SetLogger(out io.Writer) {
	//w.logger.SetOutput(out)
}

// Start initiates worker to start running the jobs.
//
// All assigned jobs will be run concurrently.
func (w *Worker) Start() {
	w.logger.Infof("running %d jobs gracefully...", len(w.jobs))

	ctx := context.Background()

	// Run registered jobs gracefully on start to
	// prevents initial burst of server load.
	go func() {
		for _, jj := range w.jobs {
			w.runner(ctx, jj)
		}
	}()

	// Handles job queueing and termination.
	for {
		select {
		// Job queue is now closed and will not run jobs anymore.
		// Queued jobs will be terminated.
		case <-w.quit:
			w.logger.Warnf("IGNO job:%s", <-w.queue)
			return
		case job := <-w.queue:
			go w.runner(ctx, job)
		}
	}
}

// AddJob registers a new job while the worker is running.
func (w *Worker) AddJob(j Job) {
	w.jobs = append(w.jobs, j)
	w.queueJob(j)
}

// runner process the job and will re-queue them when recurring job.
func (w *Worker) runner(ctx context.Context, task Job) {
	w.logger.Infof("RUNN job:%s", task)
	w.wg.Add(1)

	if err := task.Run(ctx); err != nil {
		w.logger.Errorf("ERRO job:%s - %s", task, err)
	}
	w.logger.Infof("DONE job:%s", task)
	w.wg.Done()

	// Determines if the job is run-once by interval value is zero.
	rest := task.Interval()
	if rest == 0 {
		return
	}
	// Job that has non-zero interval value means its a recurring job
	// and will be re-queued after its rest duration.
	w.logger.Infof("REST job:%s will re-queue in %s", task, rest)
	time.Sleep(rest)
	w.queueJob(task)
}

func (w *Worker) queueJob(j Job) {
	w.logger.Printf("QUED job:%s", j)
	w.queue <- j
}

// Stop will stop accepting job and wait for processing job to finish.
func (w *Worker) Stop() error {
	w.logger.Infof("stopping and waiting for jobs finish...")
	w.quit <- struct{}{}
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
