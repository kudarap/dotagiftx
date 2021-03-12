package worker

import (
	"context"
	"time"
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
