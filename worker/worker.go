package worker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

const defaultJobInterval = time.Second * 5

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
	w.logger("running...")

	go func() {
		for id, _ := range w.jobs {
			w.queueJob(id)
		}
	}()

	ctx := context.Background()

	for {
		select {
		case <-w.quit:
			return nil
		case id := <-w.queue:
			go w.runner(ctx, id)
		}
	}
}

func (w *Worker) runner(ctx context.Context, id JobID) {
	w.wg.Add(1)

	w.logger(fmt.Sprintf("job:%s recv", id))
	job := w.jobs[id]
	if err := job.Run(ctx); err != nil {
		w.logger(fmt.Sprintf("job:%s error! %s", id, err))
	}

	w.logger(fmt.Sprintf("job:%s done!", id))
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
	w.logger(fmt.Sprintf("job:%s queued", id))
}

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
