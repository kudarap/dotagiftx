package rethink

import (
	"time"

	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/gokit/log"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

// NewStats creates new instance of market data store.
func NewStats(c *Client, lg log.Logger) dgx.StatsStorage {
	return &statsStorage{c, lg}
}

type statsStorage struct {
	db     *Client
	logger log.Logger
}

func (s *statsStorage) CountUserMarketStatus(userID string) (*dgx.MarketStatusCount, error) {
	var benchStart time.Time

	baseQuery := r.Table(tableMarket).GetAllByIndex(marketFieldUserID, userID)

	var marketResult []struct {
		Group     dgx.MarketStatus `db:"group"`
		Reduction int              `db:"reduction"`
	}

	benchStart = time.Now()
	if err := s.db.list(baseQuery.
		Filter(dgx.Market{Type: dgx.MarketTypeAsk}).
		Group(marketFieldStatus).Count(), &marketResult); err != nil {
		return nil, err
	}
	mktMap := map[dgx.MarketStatus]int{}
	for _, rr := range marketResult {
		mktMap[rr.Group] = rr.Reduction
	}
	marketStats := &dgx.MarketStatusCount{
		Pending:      mktMap[dgx.MarketStatusPending],
		Live:         mktMap[dgx.MarketStatusLive],
		Sold:         mktMap[dgx.MarketStatusSold],
		Reserved:     mktMap[dgx.MarketStatusReserved],
		Removed:      mktMap[dgx.MarketStatusRemoved],
		Cancelled:    mktMap[dgx.MarketStatusCancelled],
		BidCompleted: mktMap[dgx.MarketStatusBidCompleted],
	}
	s.logger.Println("rethink/stats count ask", time.Now().Sub(benchStart))

	benchStart = time.Now()
	if err := s.db.list(baseQuery.
		HasFields(marketFieldResell).
		Filter(dgx.Market{Type: dgx.MarketTypeAsk}).
		Group(marketFieldStatus).Count(), &marketResult); err != nil {
		return nil, err
	}
	resellMap := map[dgx.MarketStatus]int{}
	for _, rr := range marketResult {
		resellMap[rr.Group] = rr.Reduction
	}
	marketStats.ResellLive = resellMap[dgx.MarketStatusLive]
	marketStats.ResellSold = resellMap[dgx.MarketStatusSold]
	marketStats.ResellReserved = resellMap[dgx.MarketStatusReserved]
	marketStats.ResellRemoved = resellMap[dgx.MarketStatusRemoved]
	marketStats.ResellCancelled = resellMap[dgx.MarketStatusCancelled]
	s.logger.Println("rethink/stats count resell", time.Now().Sub(benchStart))

	// Count market bid stats
	benchStart = time.Now()
	if err := s.db.list(baseQuery.
		Filter(dgx.Market{Type: dgx.MarketTypeBid}).
		Group(marketFieldStatus).Count(), &marketResult); err != nil {
		return nil, err
	}
	mktMap = map[dgx.MarketStatus]int{}
	for _, rr := range marketResult {
		mktMap[rr.Group] = rr.Reduction
	}
	marketStats.BidLive = mktMap[dgx.MarketStatusLive]
	marketStats.BidCompleted = mktMap[dgx.MarketStatusBidCompleted]
	s.logger.Println("rethink/stats count bid", time.Now().Sub(benchStart))

	// Count delivery stats
	benchStart = time.Now()
	var deliveryResult []struct {
		Group     dgx.DeliveryStatus `db:"group"`
		Reduction int                `db:"reduction"`
	}
	if err := s.db.list(baseQuery.Group(marketFieldDeliveryStatus).Count(), &deliveryResult); err != nil {
		return nil, err
	}
	dlvMap := map[dgx.DeliveryStatus]int{}
	for _, rr := range deliveryResult {
		dlvMap[rr.Group] = rr.Reduction
	}
	marketStats.DeliveryNoHit = dlvMap[dgx.DeliveryStatusNoHit]
	marketStats.DeliveryNameVerified = dlvMap[dgx.DeliveryStatusNameVerified]
	marketStats.DeliverySenderVerified = dlvMap[dgx.DeliveryStatusSenderVerified]
	marketStats.DeliveryPrivate = dlvMap[dgx.DeliveryStatusPrivate]
	marketStats.DeliveryError = dlvMap[dgx.DeliveryStatusError]
	s.logger.Println("rethink/stats count dlv", time.Now().Sub(benchStart))

	// Count inventory stats
	benchStart = time.Now()
	var inventoryResult []struct {
		Group     dgx.InventoryStatus `db:"group"`
		Reduction int                 `db:"reduction"`
	}
	if err := s.db.list(baseQuery.Group(marketFieldInventoryStatus).Count(), &inventoryResult); err != nil {
		return nil, err
	}
	invMap := map[dgx.InventoryStatus]int{}
	for _, rr := range inventoryResult {
		invMap[rr.Group] = rr.Reduction
	}
	marketStats.InventoryNoHit = invMap[dgx.InventoryStatusNoHit]
	marketStats.InventoryVerified = invMap[dgx.InventoryStatusVerified]
	marketStats.InventoryPrivate = invMap[dgx.InventoryStatusPrivate]
	marketStats.InventoryError = invMap[dgx.InventoryStatusError]
	s.logger.Println("rethink/stats count inv", time.Now().Sub(benchStart))

	return marketStats, nil
}

func (s *statsStorage) CountMarketStatus(opts dgx.FindOpts) (*dgx.MarketStatusCount, error) {
	var res []struct {
		Group     dgx.MarketStatus `db:"group"`
		Reduction int              `db:"reduction"`
	}
	q := newFindOptsQuery(r.Table(tableMarket).GroupByIndex(marketFieldStatus), opts)
	if err := s.db.list(q.Filter(dgx.Market{Type: dgx.MarketTypeAsk}).Count(), &res); err != nil {
		return nil, err
	}
	mapRes := map[dgx.MarketStatus]int{}
	for _, rr := range res {
		mapRes[rr.Group] = rr.Reduction
	}

	msc := &dgx.MarketStatusCount{
		Pending:      mapRes[dgx.MarketStatusPending],
		Live:         mapRes[dgx.MarketStatusLive],
		Sold:         mapRes[dgx.MarketStatusSold],
		Reserved:     mapRes[dgx.MarketStatusReserved],
		Removed:      mapRes[dgx.MarketStatusRemoved],
		Cancelled:    mapRes[dgx.MarketStatusCancelled],
		BidCompleted: mapRes[dgx.MarketStatusBidCompleted],
	}

	// Count bid stats
	q = newFindOptsQuery(r.Table(tableMarket).GroupByIndex(marketFieldStatus), opts)
	if err := s.db.list(q.Filter(dgx.Market{Type: dgx.MarketTypeBid}).Count(), &res); err != nil {
		return nil, err
	}
	mapRes = map[dgx.MarketStatus]int{}
	for _, rr := range res {
		mapRes[rr.Group] = rr.Reduction
	}
	msc.BidLive = mapRes[dgx.MarketStatusLive]
	msc.BidCompleted = mapRes[dgx.MarketStatusBidCompleted]

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

func (s *statsStorage) CountDeliveryStatus(o dgx.FindOpts) (*dgx.MarketStatusCount, error) {
	var res []struct {
		Group     dgx.DeliveryStatus `db:"group"`
		Reduction int                `db:"reduction"`
	}
	q := newFindOptsQuery(r.Table(tableMarket).GroupByIndex(marketFieldDeliveryStatus), o)
	if err := s.db.list(q.Count(), &res); err != nil {
		return nil, err
	}
	dlvMap := map[dgx.DeliveryStatus]int{}
	for _, rr := range res {
		dlvMap[rr.Group] = rr.Reduction
	}
	msc := &dgx.MarketStatusCount{
		DeliveryNoHit:          dlvMap[dgx.DeliveryStatusNoHit],
		DeliveryNameVerified:   dlvMap[dgx.DeliveryStatusNameVerified],
		DeliverySenderVerified: dlvMap[dgx.DeliveryStatusSenderVerified],
		DeliveryPrivate:        dlvMap[dgx.DeliveryStatusPrivate],
		DeliveryError:          dlvMap[dgx.DeliveryStatusError],
	}

	return msc, nil
}

func (s *statsStorage) CountInventoryStatus(o dgx.FindOpts) (*dgx.MarketStatusCount, error) {
	var res []struct {
		Group     dgx.InventoryStatus `db:"group"`
		Reduction int                 `db:"reduction"`
	}
	q := newFindOptsQuery(r.Table(tableMarket).GroupByIndex(marketFieldInventoryStatus), o)
	if err := s.db.list(q.Count(), &res); err != nil {
		return nil, err
	}
	mapRes := map[dgx.InventoryStatus]int{}
	for _, rr := range res {
		mapRes[rr.Group] = rr.Reduction
	}

	msc := &dgx.MarketStatusCount{
		InventoryNoHit:    mapRes[dgx.InventoryStatusNoHit],
		InventoryVerified: mapRes[dgx.InventoryStatusVerified],
		InventoryPrivate:  mapRes[dgx.InventoryStatusPrivate],
		InventoryError:    mapRes[dgx.InventoryStatusError],
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
func (s *statsStorage) GraphMarketSales(o dgx.FindOpts) ([]dgx.MarketSalesGraph, error) {
	o.IndexKey = marketFieldItemID
	q := newFindOptsQuery(r.Table(tableMarket), o).Filter(func(t r.Term) r.Term {
		f := t.Field(marketFieldStatus)
		return f.Eq(dgx.MarketStatusReserved).Or(f.Eq(dgx.MarketStatusSold))
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

	var msg []dgx.MarketSalesGraph
	if err := s.db.list(q, &msg); err != nil {
		return nil, err
	}
	return msg, nil
}
