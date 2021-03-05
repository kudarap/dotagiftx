package core

import "time"

type (
	// MarketStatusCount represents total number of records per status.
	MarketStatusCount struct {
		Pending      int `json:"pending"       db:"pending"`
		Live         int `json:"live"          db:"live"`
		Reserved     int `json:"reserved"      db:"reserved"`
		Sold         int `json:"sold"          db:"sold"`
		Removed      int `json:"removed"       db:"removed"`
		Cancelled    int `json:"cancelled"     db:"cancelled"`
		BidCompleted int `json:"bid_completed" db:"bid_completed"`
	}

	MarketSaleSummary struct {
		LastSalePrice     float64
		LastSaleDate      *time.Time
		LastReservedPrice float64
		LastReservedDate  *time.Time
	}

	MarketSalesGraph struct {
		Date  *time.Time `json:"date"  db:"date"`
		Avg   float64    `json:"avg"   db:"avg"`
		Count int        `json:"count" db:"count"`
	}

	// UserStats represents total users stats.
	UserStats struct {
		TotalUsersCount      int
		NewUsersThisMonth    int
		ActiveUsersLastMonth int
	}

	SalesStatus struct {
		TotalSaleValue      int
		TotalPotentialSales int
		TotalSales          int
	}

	// StatsService provides access to stats service.
	StatsService interface {
		//CountTotalMarketStatus() (*MarketStatusCount, error)
		//CountUserMarketStatus(userID string) (*MarketStatusCount, error)
		CountMarketStatus(opts FindOpts) (*MarketStatusCount, error)

		GraphMarketSales(opts FindOpts) ([]MarketSalesGraph, error)
	}

	StatsStorage StatsService
)
