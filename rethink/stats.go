package rethink

import (
	"context"
	"log/slog"
	"time"

	dotagiftx2 "github.com/kudarap/dotagiftx/dotagiftx"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

// NewStats creates new instance of market data store.
func NewStats(c *Client, lg *slog.Logger) *StatsRepository {
	return &StatsRepository{c, lg}
}

type StatsRepository struct {
	db     *Client
	logger *slog.Logger
}

func (s *StatsRepository) CountUserMarketStatus(ctx context.Context, userID string) (*dotagiftx2.MarketStatusCount, error) {
	var benchStart time.Time

	baseQuery := r.Table(tableMarket).GetAllByIndex(marketFieldUserID, userID)

	var marketResult []struct {
		Group     dotagiftx2.MarketStatus `db:"group"`
		Reduction int                     `db:"reduction"`
	}

	benchStart = time.Now()
	if err := s.db.list(ctx, baseQuery.
		Filter(dotagiftx2.Market{Type: dotagiftx2.MarketTypeAsk}).
		Group(marketFieldStatus).Count(), &marketResult); err != nil {
		return nil, err
	}
	mktMap := map[dotagiftx2.MarketStatus]int{}
	for _, rr := range marketResult {
		mktMap[rr.Group] = rr.Reduction
	}
	marketStats := &dotagiftx2.MarketStatusCount{
		Pending:      mktMap[dotagiftx2.MarketStatusPending],
		Live:         mktMap[dotagiftx2.MarketStatusLive],
		Sold:         mktMap[dotagiftx2.MarketStatusSold],
		Reserved:     mktMap[dotagiftx2.MarketStatusReserved],
		Removed:      mktMap[dotagiftx2.MarketStatusRemoved],
		Cancelled:    mktMap[dotagiftx2.MarketStatusCancelled],
		BidCompleted: mktMap[dotagiftx2.MarketStatusBidCompleted],
	}
	s.logger.InfoContext(ctx, "rethink/stats count ask", "elapsed", time.Since(benchStart))

	benchStart = time.Now()
	if err := s.db.list(ctx, baseQuery.
		HasFields(marketFieldResell).
		Filter(dotagiftx2.Market{Type: dotagiftx2.MarketTypeAsk}).
		Group(marketFieldStatus).Count(), &marketResult); err != nil {
		return nil, err
	}
	resellMap := map[dotagiftx2.MarketStatus]int{}
	for _, rr := range marketResult {
		resellMap[rr.Group] = rr.Reduction
	}
	marketStats.ResellLive = resellMap[dotagiftx2.MarketStatusLive]
	marketStats.ResellSold = resellMap[dotagiftx2.MarketStatusSold]
	marketStats.ResellReserved = resellMap[dotagiftx2.MarketStatusReserved]
	marketStats.ResellRemoved = resellMap[dotagiftx2.MarketStatusRemoved]
	marketStats.ResellCancelled = resellMap[dotagiftx2.MarketStatusCancelled]
	s.logger.InfoContext(ctx, "rethink/stats count resell", "elapsed", time.Since(benchStart))

	// Count market bid stats
	benchStart = time.Now()
	if err := s.db.list(ctx, baseQuery.
		Filter(dotagiftx2.Market{Type: dotagiftx2.MarketTypeBid}).
		Group(marketFieldStatus).Count(), &marketResult); err != nil {
		return nil, err
	}
	mktMap = map[dotagiftx2.MarketStatus]int{}
	for _, rr := range marketResult {
		mktMap[rr.Group] = rr.Reduction
	}
	marketStats.BidLive = mktMap[dotagiftx2.MarketStatusLive]
	marketStats.BidCompleted = mktMap[dotagiftx2.MarketStatusBidCompleted]
	s.logger.InfoContext(ctx, "rethink/stats count bid", "elapsed", time.Since(benchStart))

	// Count delivery stats
	benchStart = time.Now()
	var deliveryResult []struct {
		Group     dotagiftx2.DeliveryStatus `db:"group"`
		Reduction int                       `db:"reduction"`
	}
	if err := s.db.list(ctx, baseQuery.Group(marketFieldDeliveryStatus).Count(), &deliveryResult); err != nil {
		return nil, err
	}
	dlvMap := map[dotagiftx2.DeliveryStatus]int{}
	for _, rr := range deliveryResult {
		dlvMap[rr.Group] = rr.Reduction
	}
	marketStats.DeliveryNoHit = dlvMap[dotagiftx2.DeliveryStatusNoHit]
	marketStats.DeliveryNameVerified = dlvMap[dotagiftx2.DeliveryStatusNameVerified]
	marketStats.DeliverySenderVerified = dlvMap[dotagiftx2.DeliveryStatusSenderVerified]
	marketStats.DeliveryPrivate = dlvMap[dotagiftx2.DeliveryStatusPrivate]
	marketStats.DeliveryError = dlvMap[dotagiftx2.DeliveryStatusError]
	s.logger.InfoContext(ctx, "rethink/stats count dlv", "elapsed", time.Since(benchStart))

	// Count inventory stats
	benchStart = time.Now()
	var inventoryResult []struct {
		Group     dotagiftx2.InventoryStatus `db:"group"`
		Reduction int                        `db:"reduction"`
	}
	if err := s.db.list(ctx, baseQuery.Group(marketFieldInventoryStatus).Count(), &inventoryResult); err != nil {
		return nil, err
	}
	invMap := map[dotagiftx2.InventoryStatus]int{}
	for _, rr := range inventoryResult {
		invMap[rr.Group] = rr.Reduction
	}
	marketStats.InventoryNoHit = invMap[dotagiftx2.InventoryStatusNoHit]
	marketStats.InventoryVerified = invMap[dotagiftx2.InventoryStatusVerified]
	marketStats.InventoryPrivate = invMap[dotagiftx2.InventoryStatusPrivate]
	marketStats.InventoryError = invMap[dotagiftx2.InventoryStatusError]
	s.logger.InfoContext(ctx, "rethink/stats count inv", "elapsed", time.Since(benchStart))

	return marketStats, nil
}

func (s *StatsRepository) CountUserMarketStatusBySteamID(ctx context.Context, steamID string) (*dotagiftx2.MarketStatusCount, error) {
	var user dotagiftx2.User
	if err := s.db.one(ctx, r.Table(tableUser).GetAllByIndex(userFieldSteamID, steamID), &user); err != nil {
		return nil, err
	}
	return s.CountUserMarketStatus(ctx, user.ID)
}

// CountMarketStatus returns market status counts.
// TODO: optimize query because it's too slow around ~3000ms
func (s *StatsRepository) CountMarketStatus(ctx context.Context, opts dotagiftx2.FindOpts) (*dotagiftx2.MarketStatusCount, error) {
	var result []struct {
		Status dotagiftx2.MarketStatus `db:"group"`
		Count  int                     `db:"count"`
	}
	q := newFindOptsQuery(r.Table(tableMarket).GroupByIndex(marketFieldStatus), opts)
	if err := s.db.list(ctx, q.Filter(dotagiftx2.Market{Type: dotagiftx2.MarketTypeAsk}).Count(), &result); err != nil {
		return nil, err
	}
	mapResult := map[dotagiftx2.MarketStatus]int{}
	for _, v := range result {
		mapResult[v.Status] = v.Count
	}

	marketStats := &dotagiftx2.MarketStatusCount{
		Pending:      mapResult[dotagiftx2.MarketStatusPending],
		Live:         mapResult[dotagiftx2.MarketStatusLive],
		Sold:         mapResult[dotagiftx2.MarketStatusSold],
		Reserved:     mapResult[dotagiftx2.MarketStatusReserved],
		Removed:      mapResult[dotagiftx2.MarketStatusRemoved],
		Cancelled:    mapResult[dotagiftx2.MarketStatusCancelled],
		BidCompleted: mapResult[dotagiftx2.MarketStatusBidCompleted],
	}

	// Count bid stats
	q = newFindOptsQuery(r.Table(tableMarket).GroupByIndex(marketFieldStatus), opts)
	if err := s.db.list(ctx, q.Filter(dotagiftx2.Market{Type: dotagiftx2.MarketTypeBid}).Count(), &result); err != nil {
		return nil, err
	}
	mapResult = map[dotagiftx2.MarketStatus]int{}
	for _, v := range result {
		mapResult[v.Status] = v.Count
	}
	marketStats.BidLive = mapResult[dotagiftx2.MarketStatusLive]
	marketStats.BidCompleted = mapResult[dotagiftx2.MarketStatusBidCompleted]

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

func (s *StatsRepository) countDeliveryStatus(ctx context.Context, o dotagiftx2.FindOpts) (*dotagiftx2.MarketStatusCount, error) {
	var res []struct {
		Status dotagiftx2.DeliveryStatus `db:"group"`
		Count  int                       `db:"reduction"`
	}
	q := newFindOptsQuery(r.Table(tableMarket).GroupByIndex(marketFieldDeliveryStatus), o)
	if err := s.db.list(ctx, q.Count(), &res); err != nil {
		return nil, err
	}
	dlvMap := map[dotagiftx2.DeliveryStatus]int{}
	for _, v := range res {
		dlvMap[v.Status] = v.Count
	}
	msc := &dotagiftx2.MarketStatusCount{
		DeliveryNoHit:          dlvMap[dotagiftx2.DeliveryStatusNoHit],
		DeliveryNameVerified:   dlvMap[dotagiftx2.DeliveryStatusNameVerified],
		DeliverySenderVerified: dlvMap[dotagiftx2.DeliveryStatusSenderVerified],
		DeliveryPrivate:        dlvMap[dotagiftx2.DeliveryStatusPrivate],
		DeliveryError:          dlvMap[dotagiftx2.DeliveryStatusError],
	}

	return msc, nil
}

func (s *StatsRepository) countInventoryStatus(ctx context.Context, o dotagiftx2.FindOpts) (*dotagiftx2.MarketStatusCount, error) {
	var res []struct {
		Status dotagiftx2.InventoryStatus `db:"group"`
		Count  int                        `db:"reduction"`
	}
	q := newFindOptsQuery(r.Table(tableMarket).GroupByIndex(marketFieldInventoryStatus), o)
	if err := s.db.list(ctx, q.Count(), &res); err != nil {
		return nil, err
	}
	mapRes := map[dotagiftx2.InventoryStatus]int{}
	for _, rr := range res {
		mapRes[rr.Status] = rr.Count
	}

	msc := &dotagiftx2.MarketStatusCount{
		InventoryNoHit:    mapRes[dotagiftx2.InventoryStatusNoHit],
		InventoryVerified: mapRes[dotagiftx2.InventoryStatusVerified],
		InventoryPrivate:  mapRes[dotagiftx2.InventoryStatusPrivate],
		InventoryError:    mapRes[dotagiftx2.InventoryStatusError],
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
func (s *StatsRepository) GraphMarketSales(ctx context.Context, o dotagiftx2.FindOpts) ([]dotagiftx2.MarketSalesGraph, error) {
	o.IndexKey = marketFieldItemID
	q := newFindOptsQuery(r.Table(tableMarket), o).Filter(func(t r.Term) r.Term {
		f := t.Field(marketFieldStatus)
		return f.Eq(dotagiftx2.MarketStatusReserved).Or(f.Eq(dotagiftx2.MarketStatusSold))
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

	var msg []dotagiftx2.MarketSalesGraph
	if err := s.db.list(ctx, q, &msg); err != nil {
		return nil, err
	}
	return msg, nil
}
