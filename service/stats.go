package service

import (
	"github.com/kudarap/dotagiftx"
)

// NewStats returns new Stats service.
func NewStats(ss dotagiftx.StatsStorage, ts dotagiftx.TrackStorage) dotagiftx.StatsService {
	return &statsService{ss, ts}
}

type statsService struct {
	statsStg dotagiftx.StatsStorage
	trackStg dotagiftx.TrackStorage
}

func (s *statsService) CountMarketStatus(opts dotagiftx.FindOpts) (*dotagiftx.MarketStatusCount, error) {
	return s.statsStg.CountMarketStatus(opts)
}

func (s *statsService) CountTotalMarketStatus() (*dotagiftx.MarketStatusCount, error) {
	return s.statsStg.CountMarketStatus(dotagiftx.FindOpts{})
}

func (s *statsService) CountUserMarketStatus(userID string) (*dotagiftx.MarketStatusCount, error) {
	return s.statsStg.CountUserMarketStatus(userID)
}

func (s *statsService) GraphMarketSales(opts dotagiftx.FindOpts) ([]dotagiftx.MarketSalesGraph, error) {
	return s.statsStg.GraphMarketSales(opts)
}

func (s *statsService) TopKeywords() ([]dotagiftx.SearchKeywordScore, error) {
	return s.trackStg.TopKeywords()
}
