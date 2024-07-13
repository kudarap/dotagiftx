package jobs

import (
    "context"
    "time"
	
	"github.com/kudarap/dotagiftx/gokit/log"
)

type ExpiringSubscription struct {
    name     string
    interval time.Duration
	logger log.Logger
}

func (s *ExpiringSubscription) Run(ctx context.Context) error {
	bs := time.Now()
	defer func() {
		s.logger.Println("RECHECK INVENTORY BENCHMARK TIME", time.Since(bs))
	}()
	
	

	return nil
}

func (s *ExpiringSubscription) String() string {
	return s.name
}

func (s *ExpiringSubscription) Interval() time.Duration {
	return s.interval
}