package service

import "github.com/kudarap/dotagiftx/core"

// NewStats returns new Stats service.
func NewStats(ss core.StatsStorage) core.StatsService {
	return &statsService{ss}
}

type statsService struct {
	statsStg core.StatsStorage
}

func (s *statsService) CountMarketStatus(opts core.FindOpts) (*core.MarketStatusCount, error) {
	return s.statsStg.CountMarketStatus(opts)
}

func (s *statsService) CountTotalMarketStatus() (*core.MarketStatusCount, error) {
	return s.statsStg.CountMarketStatus(core.FindOpts{})
}

func (s *statsService) CountUserMarketStatus(userID string) (*core.MarketStatusCount, error) {
	return s.statsStg.CountMarketStatus(core.FindOpts{
		Filter: core.Market{UserID: userID},
	})
}
