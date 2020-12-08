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

/*
productionDB.table('market')
  .filter(r.row('status').eq(300).or(r.row('status').eq(400)))
  .group([
    r.row('updated_at').year(),
    r.row('updated_at').month(),
    r.row('updated_at').day(),
    r.row('updated_at').timezone()])
  .getField('price').ungroup()
  .map(function (doc) {
    return {
      date: r.time(doc('group').nth(0), doc('group').nth(1), doc('group').nth(2), doc('group').nth(3)),
      count: doc('reduction').count(),
      avg: doc('reduction').avg()
    }
  })
*/
func (s *statsStorage) GraphMarketSales(o core.FindOpts) ([]core.MarketSalesGraph, error) {
	q := newFindOptsQuery(r.Table(tableMarket), o).Filter(func(t r.Term) r.Term {
		f := t.Field(marketFieldStatus)
		return f.Eq(core.MarketStatusReserved).Or(f.Eq(core.MarketStatusSold))
	}).Group(func(t r.Term) []r.Term {
		f := t.Field(marketFieldUpdatedAt)
		return []r.Term{
			f.Year(),
			f.Month(),
			f.Day(),
			f.Timezone(),
		}
	}).Field(marketFieldPrice).Ungroup().Map(func(doc r.Term) interface{} {
		fg := doc.Field("group")
		fr := doc.Field("reduction")
		return map[string]interface{}{
			"date":  r.Time(fg.Nth(0), fg.Nth(1), fg.Nth(2), fg.Nth(3)),
			"count": fr.Count(),
			"avg":   fr.Avg(),
		}
	})

	var msg []core.MarketSalesGraph
	if err := s.db.list(q, &msg); err != nil {
		return nil, err
	}
	return msg, nil
}
