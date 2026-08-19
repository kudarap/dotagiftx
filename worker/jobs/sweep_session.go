package jobs

import (
	"context"
	"log/slog"
	"time"
)

type SweepSession struct {
	service sessionCleaner

	name     string
	interval time.Duration
	logger   *slog.Logger
}

func NewSweepSession(cleaner sessionCleaner, lg *slog.Logger) *SweepSession {
	return &SweepSession{
		service:  cleaner,
		name:     "job_sweep_sessions",
		interval: time.Hour * 24,
		logger:   lg,
	}
}

func (c *SweepSession) String() string { return c.name }

func (c *SweepSession) Interval() time.Duration { return c.interval }

func (c *SweepSession) Run(ctx context.Context) error {
	c.logger.Info("cleaning expired session")
	if err := c.service.PurgeExpiredSessions(ctx, time.Now()); err != nil {
		return err
	}
	c.logger.Info("expired auth session cleaned!")
	return nil
}

type sessionCleaner interface {
	PurgeExpiredSessions(ctx context.Context, from time.Time) error
}
