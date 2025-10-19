package rethink

import (
	"github.com/kudarap/dotagiftx"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

// CountMarketStatusV2 manually manages indexing for performance reasons.
func (s *statsStorage) CountMarketStatusV2(opts dotagiftx.FindOpts) (*dotagiftx.MarketStatusCount, error) {
	var result []struct {
		Status dotagiftx.MarketStatus `db:"group"`
		Count  int                    `db:"count"`
	}

	q := newFindOptsQuery(r.Table(tableMarket), opts).Group(marketFieldStatus).Count()
	if err := s.db.list(q, &result); err != nil {
		return nil, err
	}
	mapResult := map[dotagiftx.MarketStatus]int{}
	for _, v := range result {
		mapResult[v.Status] = v.Count
	}

	marketStats := &dotagiftx.MarketStatusCount{
		// sells stats
		Pending:   mapResult[dotagiftx.MarketStatusPending],
		Live:      mapResult[dotagiftx.MarketStatusLive],
		Sold:      mapResult[dotagiftx.MarketStatusSold],
		Reserved:  mapResult[dotagiftx.MarketStatusReserved],
		Removed:   mapResult[dotagiftx.MarketStatusRemoved],
		Cancelled: mapResult[dotagiftx.MarketStatusCancelled],
		// buys stats
		BidLive:      mapResult[dotagiftx.MarketStatusLive],
		BidCompleted: mapResult[dotagiftx.MarketStatusBidCompleted],
	}

	deliveryStats, err := s.countDeliveryStatusV2(opts)
	if err != nil {
		return nil, err
	}
	marketStats.DeliveryNoHit = deliveryStats.DeliveryNoHit
	marketStats.DeliveryNameVerified = deliveryStats.DeliveryNameVerified
	marketStats.DeliverySenderVerified = deliveryStats.DeliverySenderVerified
	marketStats.DeliveryPrivate = deliveryStats.DeliveryPrivate
	marketStats.DeliveryError = deliveryStats.DeliveryError

	inventoryStats, err := s.countInventoryStatusV2(opts)
	if err != nil {
		return nil, err
	}
	marketStats.InventoryNoHit = inventoryStats.InventoryNoHit
	marketStats.InventoryVerified = inventoryStats.InventoryVerified
	marketStats.InventoryPrivate = inventoryStats.InventoryPrivate
	marketStats.InventoryError = inventoryStats.InventoryError

	return marketStats, nil
}

func (s *statsStorage) countDeliveryStatusV2(o dotagiftx.FindOpts) (*dotagiftx.MarketStatusCount, error) {
	var result []struct {
		Status dotagiftx.DeliveryStatus `db:"group"`
		Count  int                      `db:"reduction"`
	}

	q := newFindOptsQuery(r.Table(tableMarket), o).Group(marketFieldDeliveryStatus).Count()
	if err := s.db.list(q, &result); err != nil {
		return nil, err
	}

	stats := map[dotagiftx.DeliveryStatus]int{}
	for _, v := range result {
		stats[v.Status] = v.Count
	}

	return &dotagiftx.MarketStatusCount{
		DeliveryNoHit:          stats[dotagiftx.DeliveryStatusNoHit],
		DeliveryNameVerified:   stats[dotagiftx.DeliveryStatusNameVerified],
		DeliverySenderVerified: stats[dotagiftx.DeliveryStatusSenderVerified],
		DeliveryPrivate:        stats[dotagiftx.DeliveryStatusPrivate],
		DeliveryError:          stats[dotagiftx.DeliveryStatusError],
	}, nil
}

func (s *statsStorage) countInventoryStatusV2(o dotagiftx.FindOpts) (*dotagiftx.MarketStatusCount, error) {
	var result []struct {
		Status dotagiftx.InventoryStatus `db:"group"`
		Count  int                       `db:"reduction"`
	}

	q := newFindOptsQuery(r.Table(tableMarket), o).Group(marketFieldInventoryStatus).Count()
	if err := s.db.list(q, &result); err != nil {
		return nil, err
	}

	stats := map[dotagiftx.InventoryStatus]int{}
	for _, rr := range result {
		stats[rr.Status] = rr.Count
	}

	return &dotagiftx.MarketStatusCount{
		InventoryNoHit:    stats[dotagiftx.InventoryStatusNoHit],
		InventoryVerified: stats[dotagiftx.InventoryStatusVerified],
		InventoryPrivate:  stats[dotagiftx.InventoryStatusPrivate],
		InventoryError:    stats[dotagiftx.InventoryStatusError],
	}, nil
}
