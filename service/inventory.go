package service

import (
	"context"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/errors"
)

// NewInventory returns new Inventory service.
func NewInventory(rs core.InventoryStorage, ms core.MarketStorage) core.InventoryService {
	return &InventoryService{rs, ms}
}

type InventoryService struct {
	InventoryStg core.InventoryStorage
	marketStg    core.MarketStorage
}

func (s *InventoryService) Inventories(opts core.FindOpts) ([]core.Inventory, *core.FindMetadata, error) {
	res, err := s.InventoryStg.Find(opts)
	if err != nil {
		return nil, nil, err
	}

	if !opts.WithMeta {
		return res, nil, err
	}

	// Get result and total count for metadata.
	tc, err := s.InventoryStg.Count(opts)
	if err != nil {
		return nil, nil, err
	}

	return res, &core.FindMetadata{
		ResultCount: len(res),
		TotalCount:  tc,
	}, nil
}

func (s *InventoryService) Inventory(id string) (*core.Inventory, error) {
	return s.InventoryStg.Get(id)
}

func (s *InventoryService) InventoryByMarketID(marketID string) (*core.Inventory, error) {
	res, err := s.InventoryStg.Find(core.FindOpts{Filter: &core.Inventory{
		MarketID: marketID,
	}})
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, core.InventoryErrNotFound
	}

	return &res[0], nil
}

func (s *InventoryService) Set(_ context.Context, inv *core.Inventory) error {
	if err := inv.CheckCreate(); err != nil {
		return errors.New(core.InventoryErrRequiredFields, err)
	}

	// Update market Inventory status.
	if err := s.marketStg.BaseUpdate(&core.Market{
		ID:              inv.MarketID,
		InventoryStatus: inv.Status,
	}); err != nil {
		return err
	}

	// Process bundle count.
	inv.BundleCount = inv.CountBundles()

	// Just update existing record.
	cur, _ := s.InventoryByMarketID(inv.MarketID)
	if cur != nil {
		inv.ID = cur.ID
		inv.Retries = cur.Retries + 1
		return s.InventoryStg.Update(inv)
	}

	return s.InventoryStg.Create(inv)
}
