package jobs

import (
	"context"
	"log/slog"
	"time"

	"github.com/kudarap/dotagiftx/dotagiftx"
	"github.com/kudarap/dotagiftx/verify"
)

// VerifyInventory represents an inventory verification job.
type VerifyInventory struct {
	inventorySvc inventoryService
	marketRepo   marketRepository
	source       *verify.Source
	logger       *slog.Logger
	// job settings
	name     string
	interval time.Duration
	filter   dotagiftx.Market
}

func NewVerifyInventory(
	is inventoryService,
	ms marketRepository,
	vs *verify.Source,
	lg *slog.Logger,
) *VerifyInventory {
	f := dotagiftx.Market{}
	return &VerifyInventory{
		is, ms, vs, lg,
		"verify_inventory", time.Hour * 24, f}
}

func (vi *VerifyInventory) String() string { return vi.name }

func (vi *VerifyInventory) Interval() time.Duration { return vi.interval }

func (vi *VerifyInventory) Run(ctx context.Context) error {
	bs := time.Now()
	defer func() {
		vi.logger.Info("VERIFIED INVENTORY BENCHMARK TIME", "elapsed", time.Since(bs))
	}()

	opts := dotagiftx.FindOpts{Filter: vi.filter}
	opts.IndexSorting = true
	opts.Sort = "updated_at"
	opts.Desc = true
	opts.Limit = 10
	opts.Page = 0

	for {
		res, err := vi.marketRepo.PendingInventoryStatus(ctx, opts)
		if err != nil {
			return err
		}

		for _, mkt := range res {
			start := time.Now()

			// Skip verified statuses.
			if mkt.InventoryStatus == dotagiftx.InventoryStatusVerified ||
				mkt.InventoryStatus == dotagiftx.InventoryStatusNoHit {

				// TODO! might remove items
				continue
			}

			if mkt.User == nil || mkt.Item == nil {
				vi.logger.Error("skipped process! missing data", "user", mkt.User, "item", mkt.Item)
				continue
			}

			result, err := vi.source.Inventory(ctx, mkt.User.SteamID, mkt.Item.Name)
			if err != nil {
				continue
			}

			vi.logger.Info("batch", "page", opts.Page, "steam_id", mkt.User.SteamID, "item", mkt.Item.Name, "status", result.Status)
			err = vi.inventorySvc.Set(ctx, &dotagiftx.Inventory{
				MarketID:   mkt.ID,
				Status:     result.Status,
				Assets:     result.Assets,
				VerifiedBy: result.VerifiedBy,
				ElapsedMs:  time.Since(start).Milliseconds(),
			})
			if err != nil {
				vi.logger.Error("inventory update error", "steam_id", mkt.User.SteamID, "item", mkt.Item.Name, "status", result.Status, "error", err)
			}

			time.Sleep(time.Second / 4)
		}

		// Is there more?
		if len(res) < opts.Limit {
			return nil
		}
		// opts.Page++
		// time.Sleep(time.Second * 2)
	}
}
