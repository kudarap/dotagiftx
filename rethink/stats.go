package rethink

import (
	"context"
	"log/slog"
	"time"

	"github.com/kudarap/dotagiftx"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

// NewStats creates new instance of market data store.
func NewStats(c *Client, lg *slog.Logger) dotagiftx.StatsStorage {
	return &statsStorage{c, lg}
}

type statsStorage struct {
	db     *Client
	logger *slog.Logger
}

func (s *statsStorage) CountUserMarketStatus(ctx context.Context, userID string) (*dotagiftx.MarketStatusCount, error) {
	var benchStart time.Time

	baseQuery := r.Table(tableMarket).GetAllByIndex(marketFieldUserID, userID)

	var marketResult []struct {
		Group     dotagiftx.MarketStatus `db:"group"`
		Reduction int                    `db:"reduction"`
	}

	benchStart = time.Now()
	if err := s.db.list(ctx, baseQuery.
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
	s.logger.Info("rethink/stats count ask", "elapsed", time.Since(benchStart))

	benchStart = time.Now()
	if err := s.db.list(ctx, baseQuery.
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
	s.logger.Info("rethink/stats count resell", "elapsed", time.Since(benchStart))

	// Count market bid stats
	benchStart = time.Now()
	if err := s.db.list(ctx, baseQuery.
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
	s.logger.Info("rethink/stats count bid", "elapsed", time.Since(benchStart))

	// Count delivery stats
	benchStart = time.Now()
	var deliveryResult []struct {
		Group     dotagiftx.DeliveryStatus `db:"group"`
		Reduction int                      `db:"reduction"`
	}
	if err := s.db.list(ctx, baseQuery.Group(marketFieldDeliveryStatus).Count(), &deliveryResult); err != nil {
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
	s.logger.Info("rethink/stats count dlv", "elapsed", time.Since(benchStart))

	// Count inventory stats
	benchStart = time.Now()
	var inventoryResult []struct {
		Group     dotagiftx.InventoryStatus `db:"group"`
		Reduction int                       `db:"reduction"`
	}
	if err := s.db.list(ctx, baseQuery.Group(marketFieldInventoryStatus).Count(), &inventoryResult); err != nil {
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
	s.logger.Info("rethink/stats count inv", "elapsed", time.Since(benchStart))

	return marketStats, nil
}

func (s *statsStorage) CountUserMarketStatusBySteamID(ctx context.Context, steamID string) (*dotagiftx.MarketStatusCount, error) {
	var user dotagiftx.User
	if err := s.db.one(ctx, r.Table(tableUser).GetAllByIndex(userFieldSteamID, steamID), &user); err != nil {
		return nil, err
	}
	return s.CountUserMarketStatus(ctx, user.ID)
}

// CountMarketStatus returns market status counts.
// TODO: optimize query because it's too slow around ~3000ms
func (s *statsStorage) CountMarketStatus(ctx context.Context, opts dotagiftx.FindOpts) (*dotagiftx.MarketStatusCount, error) {
	var result []struct {
		Status dotagiftx.MarketStatus `db:"group"`
		Count  int                    `db:"count"`
	}
	q := newFindOptsQuery(r.Table(tableMarket).GroupByIndex(marketFieldStatus), opts)
	if err := s.db.list(ctx, q.Filter(dotagiftx.Market{Type: dotagiftx.MarketTypeAsk}).Count(), &result); err != nil {
		return nil, err
	}
	mapResult := map[dotagiftx.MarketStatus]int{}
	for _, v := range result {
		mapResult[v.Status] = v.Count
	}

	marketStats := &dotagiftx.MarketStatusCount{
		Pending:      mapResult[dotagiftx.MarketStatusPending],
		Live:         mapResult[dotagiftx.MarketStatusLive],
		Sold:         mapResult[dotagiftx.MarketStatusSold],
		Reserved:     mapResult[dotagiftx.MarketStatusReserved],
		Removed:      mapResult[dotagiftx.MarketStatusRemoved],
		Cancelled:    mapResult[dotagiftx.MarketStatusCancelled],
		BidCompleted: mapResult[dotagiftx.MarketStatusBidCompleted],
	}

	// Count bid stats
	q = newFindOptsQuery(r.Table(tableMarket).GroupByIndex(marketFieldStatus), opts)
	if err := s.db.list(ctx, q.Filter(dotagiftx.Market{Type: dotagiftx.MarketTypeBid}).Count(), &result); err != nil {
		return nil, err
	}
	mapResult = map[dotagiftx.MarketStatus]int{}
	for _, v := range result {
		mapResult[v.Status] = v.Count
	}
	marketStats.BidLive = mapResult[dotagiftx.MarketStatusLive]
	marketStats.BidCompleted = mapResult[dotagiftx.MarketStatusBidCompleted]

	deliveryStats, err := s.countDeliveryStatus(ctx, opts)
	if err != nil {
		return nil, err
	}
	marketStats.DeliveryNoHit = deliveryStats.DeliveryNoHit
	marketStats.DeliveryNameVerified = deliveryStats.DeliveryNameVerified
	marketStats.DeliverySenderVerified = deliveryStats.DeliverySenderVerified
	marketStats.DeliveryPrivate = deliveryStats.DeliveryPrivate
	marketStats.DeliveryError = deliveryStats.DeliveryError

	inventoryStats, err := s.countInventoryStatus(ctx, opts)
	if err != nil {
		return nil, err
	}
	marketStats.InventoryNoHit = inventoryStats.InventoryNoHit
	marketStats.InventoryVerified = inventoryStats.InventoryVerified
	marketStats.InventoryPrivate = inventoryStats.InventoryPrivate
	marketStats.InventoryError = inventoryStats.InventoryError

	return marketStats, nil
}

func (s *statsStorage) countDeliveryStatus(ctx context.Context, o dotagiftx.FindOpts) (*dotagiftx.MarketStatusCount, error) {
	var res []struct {
		Status dotagiftx.DeliveryStatus `db:"group"`
		Count  int                      `db:"reduction"`
	}
	q := newFindOptsQuery(r.Table(tableMarket).GroupByIndex(marketFieldDeliveryStatus), o)
	if err := s.db.list(ctx, q.Count(), &res); err != nil {
		return nil, err
	}
	dlvMap := map[dotagiftx.DeliveryStatus]int{}
	for _, v := range res {
		dlvMap[v.Status] = v.Count
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

func (s *statsStorage) countInventoryStatus(ctx context.Context, o dotagiftx.FindOpts) (*dotagiftx.MarketStatusCount, error) {
	var res []struct {
		Status dotagiftx.InventoryStatus `db:"group"`
		Count  int                       `db:"reduction"`
	}
	q := newFindOptsQuery(r.Table(tableMarket).GroupByIndex(marketFieldInventoryStatus), o)
	if err := s.db.list(ctx, q.Count(), &res); err != nil {
		return nil, err
	}
	mapRes := map[dotagiftx.InventoryStatus]int{}
	for _, rr := range res {
		mapRes[rr.Status] = rr.Count
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
func (s *statsStorage) GraphMarketSales(ctx context.Context, o dotagiftx.FindOpts) ([]dotagiftx.MarketSalesGraph, error) {
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
	}).Field(marketFieldPrice).Ungroup().Map(func(doc r.Term) any {
		fg := doc.Field("group")
		fr := doc.Field("reduction")
		return map[string]any{
			"date":  r.Time(fg.Nth(0), fg.Nth(1), fg.Nth(2), fg.Nth(3)),
			"count": fr.Count(),
			"avg":   fr.Avg(),
		}
	})

	var msg []dotagiftx.MarketSalesGraph
	if err := s.db.list(ctx, q, &msg); err != nil {
		return nil, err
	}
	return msg, nil
}
