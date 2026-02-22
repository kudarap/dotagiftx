package rethink

import (
	"github.com/kudarap/dotagiftx"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

// CountMarketStatusV2 manually manages indexing for performance reasons.
func (s *statsStorage) CountMarketStatusV2(opts dotagiftx.FindOpts) (*dotagiftx.MarketStatusCount, error) {
	opts = dotagiftx.FindOpts{
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
	if err := s.db.list(q, &groups); err != nil {
		return nil, err
	}
	statusResult := map[dotagiftx.MarketStatus]int{}
	inventoryResult := map[dotagiftx.InventoryStatus]int{}
	deliveryResult := map[dotagiftx.DeliveryStatus]int{}
	resellResult := map[dotagiftx.MarketStatus]int{}
	for _, group := range groups {
		statusKey := dotagiftx.MarketStatus(group.Key[0])
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

		deliveryKey := dotagiftx.DeliveryStatus(group.Key[1])
		if _, ok := deliveryResult[deliveryKey]; !ok {
			deliveryResult[deliveryKey] = group.Count
		} else {
			deliveryResult[deliveryKey] += group.Count
		}

		inventoryKey := dotagiftx.InventoryStatus(group.Key[2])
		if _, ok := inventoryResult[inventoryKey]; !ok {
			inventoryResult[inventoryKey] = group.Count
		} else {
			inventoryResult[inventoryKey] += group.Count
		}

	}

	allStats := &dotagiftx.MarketStatusCount{
		// sells stats
		Pending:   statusResult[dotagiftx.MarketStatusPending],
		Live:      statusResult[dotagiftx.MarketStatusLive],
		Sold:      statusResult[dotagiftx.MarketStatusSold],
		Reserved:  statusResult[dotagiftx.MarketStatusReserved],
		Removed:   statusResult[dotagiftx.MarketStatusRemoved],
		Cancelled: statusResult[dotagiftx.MarketStatusCancelled],

		// buys stats
		BidLive:      statusResult[dotagiftx.MarketStatusLive],
		BidCompleted: statusResult[dotagiftx.MarketStatusBidCompleted],

		// delivery stats
		DeliveryNoHit:          deliveryResult[dotagiftx.DeliveryStatusNoHit],
		DeliveryNameVerified:   deliveryResult[dotagiftx.DeliveryStatusNameVerified],
		DeliverySenderVerified: deliveryResult[dotagiftx.DeliveryStatusSenderVerified],
		DeliveryPrivate:        deliveryResult[dotagiftx.DeliveryStatusPrivate],
		DeliveryError:          deliveryResult[dotagiftx.DeliveryStatusError],

		// inventory stats
		InventoryNoHit:    inventoryResult[dotagiftx.InventoryStatusNoHit],
		InventoryVerified: inventoryResult[dotagiftx.InventoryStatusVerified],
		InventoryPrivate:  inventoryResult[dotagiftx.InventoryStatusPrivate],
		InventoryError:    inventoryResult[dotagiftx.InventoryStatusError],

		// resell stats
		ResellLive:      resellResult[dotagiftx.MarketStatusLive],
		ResellSold:      resellResult[dotagiftx.MarketStatusSold],
		ResellReserved:  resellResult[dotagiftx.MarketStatusReserved],
		ResellRemoved:   resellResult[dotagiftx.MarketStatusRemoved],
		ResellCancelled: resellResult[dotagiftx.MarketStatusCancelled],
	}

	return allStats, nil
}
