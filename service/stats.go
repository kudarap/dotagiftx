package service

import (
	dgx "github.com/kudarap/dotagiftx"
)

// NewStats returns new Stats service.
func NewStats(ss dgx.StatsStorage, ts dgx.TrackStorage) dgx.StatsService {
	return &statsService{ss, ts}
}

type statsService struct {
	statsStg dgx.StatsStorage
	trackStg dgx.TrackStorage
}

func (s *statsService) CountMarketStatus(opts dgx.FindOpts) (*dgx.MarketStatusCount, error) {
	return s.statsStg.CountMarketStatus(opts)
}

func (s *statsService) CountTotalMarketStatus() (*dgx.MarketStatusCount, error) {
	return s.statsStg.CountMarketStatus(dgx.FindOpts{})
}

func (s *statsService) CountUserMarketStatus(userID string) (*dgx.MarketStatusCount, error) {
	return s.statsStg.CountUserMarketStatus(userID)
}

func (s *statsService) GraphMarketSales(opts dgx.FindOpts) ([]dgx.MarketSalesGraph, error) {
	return s.statsStg.GraphMarketSales(opts)
}

func (s *statsService) TopKeywords() ([]dgx.SearchKeywordScore, error) {
	return s.trackStg.TopKeywords()
}
