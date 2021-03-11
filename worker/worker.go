package worker

import (
	"context"
	"log"
	"sync"
	"time"
)

const (
	defaultJobInterval = time.Second * 5

	workerLogName = "worker"
)

type JobID string

type Job interface {
	ID() string
	Run(context.Context) error
	Interval() time.Duration
}

type Worker struct {
	wg    sync.WaitGroup
	quit  chan struct{}
	queue chan JobID
	jobs  map[JobID]Job
}

func New(jobs ...Job) *Worker {
	w := &Worker{}
	w.queue = make(chan JobID)
	w.quit = make(chan struct{})
	w.jobs = map[JobID]Job{}

	for _, jj := range jobs {
		w.jobs[JobID(jj.ID())] = jj
	}

	return w
}

func (w *Worker) Start() error {
	log.Println(workerLogName, "running...")

	ctx := context.Background()

	for {
		select {
		case <-w.quit:
			return nil
		case id := <-w.queue:
			w.wg.Add(1)
			go w.runner(ctx, id)
		}
	}
}

func (w *Worker) runner(ctx context.Context, id JobID) {
	log.Printf("[%s] job[%s] recv", workerLogName, id)
	job := w.jobs[id]
	if err := job.Run(ctx); err != nil {
		log.Printf("[%s] job[%s] error! %s", workerLogName, id, err)
	}

	log.Printf("[%s] job[%s] done!", workerLogName, id)
	w.wg.Done()

	// Rest before next iteration.
	d := job.Interval()
	if d == 0 {
		d = defaultJobInterval
	}
	time.Sleep(d)

	w.queueJob(id)
}

func (w *Worker) queueJob(id JobID) {
	w.queue <- id
}

func (w *Worker) Stop() error {
	log.Println(workerLogName, "stopping and waiting for jobs finish...")
	w.quit <- struct{}{}
	w.wg.Wait()
	log.Println(workerLogName, "all jobs done!")
	return nil
}
