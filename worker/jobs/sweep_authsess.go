package jobs

import (
	"context"
	"log/slog"
	"time"
)

type CleanAuthSession struct {
	service authSessCleaner

	name     string
	interval time.Duration
	logger   *slog.Logger
}

func NewSweepAuthSess(cleaner authSessCleaner, lg *slog.Logger) *CleanAuthSession {
	return &CleanAuthSession{
		service:  cleaner,
		name:     "job_sweep_auth_sess",
		interval: time.Hour * 24,
		logger:   lg,
	}
}

func (c *CleanAuthSession) String() string { return c.name }

func (c *CleanAuthSession) Interval() time.Duration { return c.interval }

func (c *CleanAuthSession) Run(ctx context.Context) error {
	c.logger.Info("cleaning expired auth session")
	if err := c.service.CleanExpiredSession(ctx, time.Now()); err != nil {
		return err
	}
	c.logger.Info("expired auth session cleaned!")
	return nil
}

type authSessCleaner interface {
	CleanExpiredSession(ctx context.Context, from time.Time) error
}
