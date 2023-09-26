package rethink

import (
	"fmt"
	"time"

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

func (s *statsStorage) CountUserMarketStatus(userID string) (*core.MarketStatusCount, error) {
	var benchStart time.Time

	baseQuery := r.Table(tableMarket).GetAllByIndex(marketFieldUserID, userID)

	var marketResult []struct {
		Group     core.MarketStatus `db:"group"`
		Reduction int               `db:"reduction"`
	}

	benchStart = time.Now()
	if err := s.db.list(baseQuery.
		Filter(core.Market{Type: core.MarketTypeAsk}).
		Group(marketFieldStatus).Count(), &marketResult); err != nil {
		return nil, err
	}
	mktMap := map[core.MarketStatus]int{}
	for _, rr := range marketResult {
		mktMap[rr.Group] = rr.Reduction
	}
	marketStats := &core.MarketStatusCount{
		Pending:      mktMap[core.MarketStatusPending],
		Live:         mktMap[core.MarketStatusLive],
		Sold:         mktMap[core.MarketStatusSold],
		Reserved:     mktMap[core.MarketStatusReserved],
		Removed:      mktMap[core.MarketStatusRemoved],
		Cancelled:    mktMap[core.MarketStatusCancelled],
		BidCompleted: mktMap[core.MarketStatusBidCompleted],
	}
	fmt.Println("rethink/stats count ask", time.Now().Sub(benchStart))

	benchStart = time.Now()
	if err := s.db.list(baseQuery.
		HasFields(marketFieldResell).
		Filter(core.Market{Type: core.MarketTypeAsk}).
		Group(marketFieldStatus).Count(), &marketResult); err != nil {
		return nil, err
	}
	resellMap := map[core.MarketStatus]int{}
	for _, rr := range marketResult {
		resellMap[rr.Group] = rr.Reduction
	}
	marketStats.ResellLive = resellMap[core.MarketStatusLive]
	marketStats.ResellSold = resellMap[core.MarketStatusSold]
	marketStats.ResellReserved = resellMap[core.MarketStatusReserved]
	marketStats.ResellRemoved = resellMap[core.MarketStatusRemoved]
	marketStats.ResellCancelled = resellMap[core.MarketStatusCancelled]
	fmt.Println("rethink/stats count resell", time.Now().Sub(benchStart))

	// Count market bid stats
	benchStart = time.Now()
	if err := s.db.list(baseQuery.
		Filter(core.Market{Type: core.MarketTypeBid}).
		Group(marketFieldStatus).Count(), &marketResult); err != nil {
		return nil, err
	}
	mktMap = map[core.MarketStatus]int{}
	for _, rr := range marketResult {
		mktMap[rr.Group] = rr.Reduction
	}
	marketStats.BidLive = mktMap[core.MarketStatusLive]
	marketStats.BidCompleted = mktMap[core.MarketStatusBidCompleted]
	fmt.Println("rethink/stats count bid", time.Now().Sub(benchStart))

	// Count delivery stats
	benchStart = time.Now()
	var deliveryResult []struct {
		Group     core.DeliveryStatus `db:"group"`
		Reduction int                 `db:"reduction"`
	}
	if err := s.db.list(baseQuery.Group(marketFieldDeliveryStatus).Count(), &deliveryResult); err != nil {
		return nil, err
	}
	dlvMap := map[core.DeliveryStatus]int{}
	for _, rr := range deliveryResult {
		dlvMap[rr.Group] = rr.Reduction
	}
	marketStats.DeliveryNoHit = dlvMap[core.DeliveryStatusNoHit]
	marketStats.DeliveryNameVerified = dlvMap[core.DeliveryStatusNameVerified]
	marketStats.DeliverySenderVerified = dlvMap[core.DeliveryStatusSenderVerified]
	marketStats.DeliveryPrivate = dlvMap[core.DeliveryStatusPrivate]
	marketStats.DeliveryError = dlvMap[core.DeliveryStatusError]
	fmt.Println("rethink/stats count dlv", time.Now().Sub(benchStart))

	// Count inventory stats
	benchStart = time.Now()
	var inventoryResult []struct {
		Group     core.InventoryStatus `db:"group"`
		Reduction int                  `db:"reduction"`
	}
	if err := s.db.list(baseQuery.Group(marketFieldInventoryStatus).Count(), &inventoryResult); err != nil {
		return nil, err
	}
	invMap := map[core.InventoryStatus]int{}
	for _, rr := range inventoryResult {
		invMap[rr.Group] = rr.Reduction
	}
	marketStats.InventoryNoHit = invMap[core.InventoryStatusNoHit]
	marketStats.InventoryVerified = invMap[core.InventoryStatusVerified]
	marketStats.InventoryPrivate = invMap[core.InventoryStatusPrivate]
	marketStats.InventoryError = invMap[core.InventoryStatusError]
	fmt.Println("rethink/stats count inv", time.Now().Sub(benchStart))

	return marketStats, nil
}

func (s *statsStorage) CountMarketStatus(opts core.FindOpts) (*core.MarketStatusCount, error) {
	var res []struct {
		Group     core.MarketStatus `db:"group"`
		Reduction int               `db:"reduction"`
	}
	q := newFindOptsQuery(r.Table(tableMarket).GroupByIndex(marketFieldStatus), opts)
	if err := s.db.list(q.Filter(core.Market{Type: core.MarketTypeAsk}).Count(), &res); err != nil {
		return nil, err
	}
	mapRes := map[core.MarketStatus]int{}
	for _, rr := range res {
		mapRes[rr.Group] = rr.Reduction
	}

	msc := &core.MarketStatusCount{
		Pending:      mapRes[core.MarketStatusPending],
		Live:         mapRes[core.MarketStatusLive],
		Sold:         mapRes[core.MarketStatusSold],
		Reserved:     mapRes[core.MarketStatusReserved],
		Removed:      mapRes[core.MarketStatusRemoved],
		Cancelled:    mapRes[core.MarketStatusCancelled],
		BidCompleted: mapRes[core.MarketStatusBidCompleted],
	}

	// Count bid stats
	q = newFindOptsQuery(r.Table(tableMarket).GroupByIndex(marketFieldStatus), opts)
	if err := s.db.list(q.Filter(core.Market{Type: core.MarketTypeBid}).Count(), &res); err != nil {
		return nil, err
	}
	mapRes = map[core.MarketStatus]int{}
	for _, rr := range res {
		mapRes[rr.Group] = rr.Reduction
	}
	msc.BidLive = mapRes[core.MarketStatusLive]
	msc.BidCompleted = mapRes[core.MarketStatusBidCompleted]

	cds, err := s.CountDeliveryStatus(opts)
	if err != nil {
		return nil, err
	}
	msc.DeliveryNoHit = cds.DeliveryNoHit
	msc.DeliveryNameVerified = cds.DeliveryNameVerified
	msc.DeliverySenderVerified = cds.DeliverySenderVerified
	msc.DeliveryPrivate = cds.DeliveryPrivate
	msc.DeliveryError = cds.DeliveryError

	cis, err := s.CountInventoryStatus(opts)
	if err != nil {
		return nil, err
	}
	msc.InventoryNoHit = cis.InventoryNoHit
	msc.InventoryVerified = cis.InventoryVerified
	msc.InventoryPrivate = cis.InventoryPrivate
	msc.InventoryError = cis.InventoryError

	return msc, nil
}

func (s *statsStorage) CountDeliveryStatus(o core.FindOpts) (*core.MarketStatusCount, error) {
	var res []struct {
		Group     core.DeliveryStatus `db:"group"`
		Reduction int                 `db:"reduction"`
	}
	q := newFindOptsQuery(r.Table(tableMarket).GroupByIndex(marketFieldDeliveryStatus), o)
	if err := s.db.list(q.Count(), &res); err != nil {
		return nil, err
	}
	dlvMap := map[core.DeliveryStatus]int{}
	for _, rr := range res {
		dlvMap[rr.Group] = rr.Reduction
	}
	msc := &core.MarketStatusCount{
		DeliveryNoHit:          dlvMap[core.DeliveryStatusNoHit],
		DeliveryNameVerified:   dlvMap[core.DeliveryStatusNameVerified],
		DeliverySenderVerified: dlvMap[core.DeliveryStatusSenderVerified],
		DeliveryPrivate:        dlvMap[core.DeliveryStatusPrivate],
		DeliveryError:          dlvMap[core.DeliveryStatusError],
	}

	return msc, nil
}

func (s *statsStorage) CountInventoryStatus(o core.FindOpts) (*core.MarketStatusCount, error) {
	var res []struct {
		Group     core.InventoryStatus `db:"group"`
		Reduction int                  `db:"reduction"`
	}
	q := newFindOptsQuery(r.Table(tableMarket).GroupByIndex(marketFieldInventoryStatus), o)
	if err := s.db.list(q.Count(), &res); err != nil {
		return nil, err
	}
	mapRes := map[core.InventoryStatus]int{}
	for _, rr := range res {
		mapRes[rr.Group] = rr.Reduction
	}

	msc := &core.MarketStatusCount{
		InventoryNoHit:    mapRes[core.InventoryStatusNoHit],
		InventoryVerified: mapRes[core.InventoryStatusVerified],
		InventoryPrivate:  mapRes[core.InventoryStatusPrivate],
		InventoryError:    mapRes[core.InventoryStatusError],
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
	o.IndexKey = marketFieldItemID
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
