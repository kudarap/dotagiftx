package rethink

import (
	"time"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/logging"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

// NewStats creates new instance of market data store.
func NewStats(c *Client, lg logging.Logger) dotagiftx.StatsStorage {
	return &statsStorage{c, lg}
}

type statsStorage struct {
	db     *Client
	logger logging.Logger
}

func (s *statsStorage) CountUserMarketStatus(userID string) (*dotagiftx.MarketStatusCount, error) {
	var benchStart time.Time

	baseQuery := r.Table(tableMarket).GetAllByIndex(marketFieldUserID, userID)

	var marketResult []struct {
		Group     dotagiftx.MarketStatus `db:"group"`
		Reduction int                    `db:"reduction"`
	}

	benchStart = time.Now()
	if err := s.db.list(baseQuery.
		Filter(dotagiftx.Market{Type: dotagiftx.MarketTypeAsk}).
		Group(marketFieldStatus).Count(), &marketResult); err != nil {
		return nil, err
	}
	mktMap := map[dotagiftx.MarketStatus]int{}
	for _, rr := range marketResult {
		mktMap[rr.Group] = rr.Reduction
	}
	marketStats := &dotagiftx.MarketStatusCount{
		Pending:      mktMap[dotagiftx.MarketStatusPending],
		Live:         mktMap[dotagiftx.MarketStatusLive],
		Sold:         mktMap[dotagiftx.MarketStatusSold],
		Reserved:     mktMap[dotagiftx.MarketStatusReserved],
		Removed:      mktMap[dotagiftx.MarketStatusRemoved],
		Cancelled:    mktMap[dotagiftx.MarketStatusCancelled],
		BidCompleted: mktMap[dotagiftx.MarketStatusBidCompleted],
	}
	s.logger.Println("rethink/stats count ask", time.Since(benchStart))

	benchStart = time.Now()
	if err := s.db.list(baseQuery.
		HasFields(marketFieldResell).
		Filter(dotagiftx.Market{Type: dotagiftx.MarketTypeAsk}).
		Group(marketFieldStatus).Count(), &marketResult); err != nil {
		return nil, err
	}
	resellMap := map[dotagiftx.MarketStatus]int{}
	for _, rr := range marketResult {
		resellMap[rr.Group] = rr.Reduction
	}
	marketStats.ResellLive = resellMap[dotagiftx.MarketStatusLive]
	marketStats.ResellSold = resellMap[dotagiftx.MarketStatusSold]
	marketStats.ResellReserved = resellMap[dotagiftx.MarketStatusReserved]
	marketStats.ResellRemoved = resellMap[dotagiftx.MarketStatusRemoved]
	marketStats.ResellCancelled = resellMap[dotagiftx.MarketStatusCancelled]
	s.logger.Println("rethink/stats count resell", time.Since(benchStart))

	// Count market bid stats
	benchStart = time.Now()
	if err := s.db.list(baseQuery.
		Filter(dotagiftx.Market{Type: dotagiftx.MarketTypeBid}).
		Group(marketFieldStatus).Count(), &marketResult); err != nil {
		return nil, err
	}
	mktMap = map[dotagiftx.MarketStatus]int{}
	for _, rr := range marketResult {
		mktMap[rr.Group] = rr.Reduction
	}
	marketStats.BidLive = mktMap[dotagiftx.MarketStatusLive]
	marketStats.BidCompleted = mktMap[dotagiftx.MarketStatusBidCompleted]
	s.logger.Println("rethink/stats count bid", time.Since(benchStart))

	// Count delivery stats
	benchStart = time.Now()
	var deliveryResult []struct {
		Group     dotagiftx.DeliveryStatus `db:"group"`
		Reduction int                      `db:"reduction"`
	}
	if err := s.db.list(baseQuery.Group(marketFieldDeliveryStatus).Count(), &deliveryResult); err != nil {
		return nil, err
	}
	dlvMap := map[dotagiftx.DeliveryStatus]int{}
	for _, rr := range deliveryResult {
		dlvMap[rr.Group] = rr.Reduction
	}
	marketStats.DeliveryNoHit = dlvMap[dotagiftx.DeliveryStatusNoHit]
	marketStats.DeliveryNameVerified = dlvMap[dotagiftx.DeliveryStatusNameVerified]
	marketStats.DeliverySenderVerified = dlvMap[dotagiftx.DeliveryStatusSenderVerified]
	marketStats.DeliveryPrivate = dlvMap[dotagiftx.DeliveryStatusPrivate]
	marketStats.DeliveryError = dlvMap[dotagiftx.DeliveryStatusError]
	s.logger.Println("rethink/stats count dlv", time.Since(benchStart))

	// Count inventory stats
	benchStart = time.Now()
	var inventoryResult []struct {
		Group     dotagiftx.InventoryStatus `db:"group"`
		Reduction int                       `db:"reduction"`
	}
	if err := s.db.list(baseQuery.Group(marketFieldInventoryStatus).Count(), &inventoryResult); err != nil {
		return nil, err
	}
	invMap := map[dotagiftx.InventoryStatus]int{}
	for _, rr := range inventoryResult {
		invMap[rr.Group] = rr.Reduction
	}
	marketStats.InventoryNoHit = invMap[dotagiftx.InventoryStatusNoHit]
	marketStats.InventoryVerified = invMap[dotagiftx.InventoryStatusVerified]
	marketStats.InventoryPrivate = invMap[dotagiftx.InventoryStatusPrivate]
	marketStats.InventoryError = invMap[dotagiftx.InventoryStatusError]
	s.logger.Println("rethink/stats count inv", time.Since(benchStart))

	return marketStats, nil
}

func (s *statsStorage) CountUserMarketStatusBySteamID(steamID string) (*dotagiftx.MarketStatusCount, error) {
	var user dotagiftx.User
	if err := s.db.one(r.Table(tableUser).GetAllByIndex(userFieldSteamID, steamID), &user); err != nil {
		return nil, err
	}
	return s.CountUserMarketStatus(user.ID)
}

// CountMarketStatus returns market status counts.
// TODO: optimize query because it's too slow around ~3000ms'
func (s *statsStorage) CountMarketStatus(opts dotagiftx.FindOpts) (*dotagiftx.MarketStatusCount, error) {
	var res []struct {
		Group     dotagiftx.MarketStatus `db:"group"`
		Reduction int                    `db:"reduction"`
	}
	q := newFindOptsQuery(r.Table(tableMarket).GroupByIndex(marketFieldStatus), opts)
	if err := s.db.list(q.Filter(dotagiftx.Market{Type: dotagiftx.MarketTypeAsk}).Count(), &res); err != nil {
		return nil, err
	}
	mapRes := map[dotagiftx.MarketStatus]int{}
	for _, rr := range res {
		mapRes[rr.Group] = rr.Reduction
	}

	msc := &dotagiftx.MarketStatusCount{
		Pending:      mapRes[dotagiftx.MarketStatusPending],
		Live:         mapRes[dotagiftx.MarketStatusLive],
		Sold:         mapRes[dotagiftx.MarketStatusSold],
		Reserved:     mapRes[dotagiftx.MarketStatusReserved],
		Removed:      mapRes[dotagiftx.MarketStatusRemoved],
		Cancelled:    mapRes[dotagiftx.MarketStatusCancelled],
		BidCompleted: mapRes[dotagiftx.MarketStatusBidCompleted],
	}

	// Count bid stats
	q = newFindOptsQuery(r.Table(tableMarket).GroupByIndex(marketFieldStatus), opts)
	if err := s.db.list(q.Filter(dotagiftx.Market{Type: dotagiftx.MarketTypeBid}).Count(), &res); err != nil {
		return nil, err
	}
	mapRes = map[dotagiftx.MarketStatus]int{}
	for _, rr := range res {
		mapRes[rr.Group] = rr.Reduction
	}
	msc.BidLive = mapRes[dotagiftx.MarketStatusLive]
	msc.BidCompleted = mapRes[dotagiftx.MarketStatusBidCompleted]

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

func (s *statsStorage) CountDeliveryStatus(o dotagiftx.FindOpts) (*dotagiftx.MarketStatusCount, error) {
	var res []struct {
		Group     dotagiftx.DeliveryStatus `db:"group"`
		Reduction int                      `db:"reduction"`
	}
	q := newFindOptsQuery(r.Table(tableMarket).GroupByIndex(marketFieldDeliveryStatus), o)
	if err := s.db.list(q.Count(), &res); err != nil {
		return nil, err
	}
	dlvMap := map[dotagiftx.DeliveryStatus]int{}
	for _, rr := range res {
		dlvMap[rr.Group] = rr.Reduction
	}
	msc := &dotagiftx.MarketStatusCount{
		DeliveryNoHit:          dlvMap[dotagiftx.DeliveryStatusNoHit],
		DeliveryNameVerified:   dlvMap[dotagiftx.DeliveryStatusNameVerified],
		DeliverySenderVerified: dlvMap[dotagiftx.DeliveryStatusSenderVerified],
		DeliveryPrivate:        dlvMap[dotagiftx.DeliveryStatusPrivate],
		DeliveryError:          dlvMap[dotagiftx.DeliveryStatusError],
	}

	return msc, nil
}

func (s *statsStorage) CountInventoryStatus(o dotagiftx.FindOpts) (*dotagiftx.MarketStatusCount, error) {
	var res []struct {
		Group     dotagiftx.InventoryStatus `db:"group"`
		Reduction int                       `db:"reduction"`
	}
	q := newFindOptsQuery(r.Table(tableMarket).GroupByIndex(marketFieldInventoryStatus), o)
	if err := s.db.list(q.Count(), &res); err != nil {
		return nil, err
	}
	mapRes := map[dotagiftx.InventoryStatus]int{}
	for _, rr := range res {
		mapRes[rr.Group] = rr.Reduction
	}

	msc := &dotagiftx.MarketStatusCount{
		InventoryNoHit:    mapRes[dotagiftx.InventoryStatusNoHit],
		InventoryVerified: mapRes[dotagiftx.InventoryStatusVerified],
		InventoryPrivate:  mapRes[dotagiftx.InventoryStatusPrivate],
		InventoryError:    mapRes[dotagiftx.InventoryStatusError],
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
func (s *statsStorage) GraphMarketSales(o dotagiftx.FindOpts) ([]dotagiftx.MarketSalesGraph, error) {
	o.IndexKey = marketFieldItemID
	q := newFindOptsQuery(r.Table(tableMarket), o).Filter(func(t r.Term) r.Term {
		f := t.Field(marketFieldStatus)
		return f.Eq(dotagiftx.MarketStatusReserved).Or(f.Eq(dotagiftx.MarketStatusSold))
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

	var msg []dotagiftx.MarketSalesGraph
	if err := s.db.list(q, &msg); err != nil {
		return nil, err
	}
	return msg, nil
}
