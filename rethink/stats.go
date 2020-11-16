package rethink

import (
	"github.com/kudarap/dotagiftx/core"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

// NewStats creates new instance of market data store.
func NewStats(c *Client) core.StatsStorage {
	return &statsStorage{c}
}

type statsStorage struct {
	db *Client
}

func (s *statsStorage) CountMarketStatus(o core.FindOpts) (*core.MarketStatusCount, error) {
	var res []struct {
		Group     core.MarketStatus `db:"group"`
		Reduction int               `db:"reduction"`
	}
	q := newFindOptsQuery(r.Table(tableMarket).GroupByIndex(marketFieldStatus), o)
	if err := s.db.list(q.Count(), &res); err != nil {
		return nil, err
	}
	mapRes := map[core.MarketStatus]int{}
	for _, rr := range res {
		mapRes[rr.Group] = rr.Reduction
	}

	msc := &core.MarketStatusCount{
		Pending:   mapRes[core.MarketStatusPending],
		Live:      mapRes[core.MarketStatusLive],
		Sold:      mapRes[core.MarketStatusSold],
		Reserved:  mapRes[core.MarketStatusReserved],
		Removed:   mapRes[core.MarketStatusRemoved],
		Cancelled: mapRes[core.MarketStatusCancelled],
	}

	return msc, nil
}
