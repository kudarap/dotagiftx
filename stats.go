package dotagiftx

import (
	"context"
	"time"
)

type (
	// MarketStatusCount represents the total number of records per status.
	MarketStatusCount struct {
		Pending   int `json:"pending"   db:"pending"`
		Live      int `json:"live"      db:"live"`
		Reserved  int `json:"reserved"  db:"reserved"`
		Sold      int `json:"sold"      db:"sold"`
		Removed   int `json:"removed"   db:"removed"`
		Cancelled int `json:"cancelled" db:"cancelled"`

		BidLive      int `json:"bid_live"      db:"bid_live"`
		BidCompleted int `json:"bid_completed" db:"bid_completed"`

		DeliveryNoHit          int `json:"delivery_no_hit"          db:"delivery_no_hit"`
		DeliveryNameVerified   int `json:"delivery_name_verified"   db:"delivery_name_verified"`
		DeliverySenderVerified int `json:"delivery_sender_verified" db:"delivery_sender_verified"`
		DeliveryPrivate        int `json:"delivery_private"         db:"delivery_private"`
		DeliveryError          int `json:"delivery_error"           db:"delivery_error"`

		InventoryNoHit    int `json:"inventory_no_hit"   db:"inventory_no_hit"`
		InventoryVerified int `json:"inventory_verified" db:"inventory_verified"`
		InventoryPrivate  int `json:"inventory_private"  db:"inventory_private"`
		InventoryError    int `json:"inventory_error"    db:"inventory_error"`

		ResellLive      int `json:"resell_live" db:"resell_live"`
		ResellReserved  int `json:"resell_reserved" db:"resell_reserved"`
		ResellSold      int `json:"resell_sold" db:"resell_sold"`
		ResellRemoved   int `json:"resell_removed" db:"resell_removed"`
		ResellCancelled int `json:"resell_cancelled" db:"resell_cancelled"`
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

	SearchKeywordScore struct {
		Keyword string `json:"keyword"`
		Score   int    `json:"score"`
	}

	statsStorage interface {
		CountMarketStatus(ctx context.Context, opts FindOpts) (*MarketStatusCount, error)
		CountMarketStatusV2(ctx context.Context, opts FindOpts) (*MarketStatusCount, error)

		GraphMarketSales(ctx context.Context, opts FindOpts) ([]MarketSalesGraph, error)

		CountUserMarketStatus(ctx context.Context, userID string) (*MarketStatusCount, error)

		CountUserMarketStatusBySteamID(ctx context.Context, partnerSteamID string) (*MarketStatusCount, error)
	}
)

// NewStatsService returns new Stats service.
func NewStatsService(ss statsStorage, ts trackStorage) *StatsService {
	return &StatsService{ss, ts}
}

type StatsService struct {
	statsStg statsStorage
	trackStg trackStorage
}

func (s *StatsService) CountUserMarketStatusBySteamID(ctx context.Context, partnerSteamID string) (*MarketStatusCount, error) {
	return s.statsStg.CountUserMarketStatusBySteamID(ctx, partnerSteamID)
}

func (s *StatsService) CountMarketStatus(ctx context.Context, opts FindOpts) (*MarketStatusCount, error) {
	return s.statsStg.CountMarketStatus(ctx, opts)
}

func (s *StatsService) CountMarketStatusV2(ctx context.Context, opts FindOpts) (*MarketStatusCount, error) {
	return s.statsStg.CountMarketStatusV2(ctx, opts)
}

func (s *StatsService) CountUserMarketStatus(ctx context.Context, userID string) (*MarketStatusCount, error) {
	return s.statsStg.CountUserMarketStatus(ctx, userID)
}

func (s *StatsService) GraphMarketSales(ctx context.Context, opts FindOpts) ([]MarketSalesGraph, error) {
	return s.statsStg.GraphMarketSales(ctx, opts)
}

func (s *StatsService) TopKeywords(ctx context.Context) ([]SearchKeywordScore, error) {
	return s.trackStg.TopKeywords(ctx)
}
