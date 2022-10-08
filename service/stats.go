package service

import "github.com/kudarap/dotagiftx/core"

// NewStats returns new Stats service.
func NewStats(ss core.StatsStorage, ts core.TrackStorage) core.StatsService {
	return &statsService{ss, ts}
}

type statsService struct {
	statsStg core.StatsStorage
	trackStg core.TrackStorage
}

func (s *statsService) CountMarketStatus(opts core.FindOpts) (*core.MarketStatusCount, error) {
	return s.statsStg.CountMarketStatus(opts)
}

func (s *statsService) CountTotalMarketStatus() (*core.MarketStatusCount, error) {
	return s.statsStg.CountMarketStatus(core.FindOpts{})
}

func (s *statsService) CountUserMarketStatus(userID string) (*core.MarketStatusCount, error) {
	return s.statsStg.CountUserMarketStatus(userID)
}

func (s *statsService) GraphMarketSales(opts core.FindOpts) ([]core.MarketSalesGraph, error) {
	return s.statsStg.GraphMarketSales(opts)
}

func (s *statsService) TopKeywords() ([]core.SearchKeywordScore, error) {
	return s.trackStg.TopKeywords()
}
