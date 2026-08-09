package rethink

import (
	"context"

	dotagiftx2 "github.com/kudarap/dotagiftx/dotagiftx"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

// CountMarketStatusV2 manually manages indexing for performance reasons.
func (s *StatsRepository) CountMarketStatusV2(ctx context.Context, opts dotagiftx2.FindOpts) (*dotagiftx2.MarketStatusCount, error) {
	opts = dotagiftx2.FindOpts{
		Filter:   opts.Filter,
		IndexKey: opts.IndexKey,
		UserID:   opts.UserID,
	}
	if opts.IndexKey == "" {
		opts.IndexKey = marketFieldType
	}

	var groups []struct {
		Key   [4]uint `db:"group"`
		Count int     `db:"reduction"`
	}

	q := newFindOptsQuery(r.Table(tableMarket), opts).MultiGroup(
		marketFieldStatus,
		marketFieldDeliveryStatus,
		marketFieldInventoryStatus,
		marketFieldResell,
	).Count()
	if err := s.db.list(ctx, q, &groups); err != nil {
		return nil, err
	}
	statusResult := map[dotagiftx2.MarketStatus]int{}
	inventoryResult := map[dotagiftx2.InventoryStatus]int{}
	deliveryResult := map[dotagiftx2.DeliveryStatus]int{}
	resellResult := map[dotagiftx2.MarketStatus]int{}
	for _, group := range groups {
		statusKey := dotagiftx2.MarketStatus(group.Key[0])
		if _, ok := statusResult[statusKey]; !ok {
			statusResult[statusKey] = group.Count
		} else {
			statusResult[statusKey] += group.Count
		}

		isResell := group.Key[3] == 1
		if isResell {
			if _, ok := resellResult[statusKey]; !ok {
				resellResult[statusKey] = group.Count
			} else {
				resellResult[statusKey] += group.Count
			}
		}

		deliveryKey := dotagiftx2.DeliveryStatus(group.Key[1])
		if _, ok := deliveryResult[deliveryKey]; !ok {
			deliveryResult[deliveryKey] = group.Count
		} else {
			deliveryResult[deliveryKey] += group.Count
		}

		inventoryKey := dotagiftx2.InventoryStatus(group.Key[2])
		if _, ok := inventoryResult[inventoryKey]; !ok {
			inventoryResult[inventoryKey] = group.Count
		} else {
			inventoryResult[inventoryKey] += group.Count
		}

	}

	allStats := &dotagiftx2.MarketStatusCount{
		// sells stats
		Pending:   statusResult[dotagiftx2.MarketStatusPending],
		Live:      statusResult[dotagiftx2.MarketStatusLive],
		Sold:      statusResult[dotagiftx2.MarketStatusSold],
		Reserved:  statusResult[dotagiftx2.MarketStatusReserved],
		Removed:   statusResult[dotagiftx2.MarketStatusRemoved],
		Cancelled: statusResult[dotagiftx2.MarketStatusCancelled],

		// buys stats
		BidLive:      statusResult[dotagiftx2.MarketStatusLive],
		BidCompleted: statusResult[dotagiftx2.MarketStatusBidCompleted],

		// delivery stats
		DeliveryNoHit:          deliveryResult[dotagiftx2.DeliveryStatusNoHit],
		DeliveryNameVerified:   deliveryResult[dotagiftx2.DeliveryStatusNameVerified],
		DeliverySenderVerified: deliveryResult[dotagiftx2.DeliveryStatusSenderVerified],
		DeliveryPrivate:        deliveryResult[dotagiftx2.DeliveryStatusPrivate],
		DeliveryError:          deliveryResult[dotagiftx2.DeliveryStatusError],

		// inventory stats
		InventoryNoHit:    inventoryResult[dotagiftx2.InventoryStatusNoHit],
		InventoryVerified: inventoryResult[dotagiftx2.InventoryStatusVerified],
		InventoryPrivate:  inventoryResult[dotagiftx2.InventoryStatusPrivate],
		InventoryError:    inventoryResult[dotagiftx2.InventoryStatusError],

		// resell stats
		ResellLive:      resellResult[dotagiftx2.MarketStatusLive],
		ResellSold:      resellResult[dotagiftx2.MarketStatusSold],
		ResellReserved:  resellResult[dotagiftx2.MarketStatusReserved],
		ResellRemoved:   resellResult[dotagiftx2.MarketStatusRemoved],
		ResellCancelled: resellResult[dotagiftx2.MarketStatusCancelled],
	}

	return allStats, nil
}
