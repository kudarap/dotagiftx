package service

import (
	"context"
	"log"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/errors"
)

// NewInventory returns new Inventory service.
func NewInventory(rs core.InventoryStorage, ms core.MarketStorage) core.InventoryService {
	return &InventoryService{rs, ms}
}

type InventoryService struct {
	inventoryStg core.InventoryStorage
	marketStg    core.MarketStorage
}

func (s *InventoryService) Inventories(opts core.FindOpts) ([]core.Inventory, *core.FindMetadata, error) {
	res, err := s.inventoryStg.Find(opts)
	if err != nil {
		return nil, nil, err
	}

	if !opts.WithMeta {
		return res, nil, err
	}

	// Get result and total count for metadata.
	tc, err := s.inventoryStg.Count(opts)
	if err != nil {
		return nil, nil, err
	}

	return res, &core.FindMetadata{
		ResultCount: len(res),
		TotalCount:  tc,
	}, nil
}

func (s *InventoryService) Inventory(id string) (*core.Inventory, error) {
	inv, err := s.inventoryStg.Get(id)
	if err != nil && err != core.InventoryErrNotFound {
		return nil, err
	}
	if inv != nil {
		return inv, nil
	}

	// If we can't find using id. lets try market ID
	return s.inventoryStg.GetByMarketID(id)
}

func (s *InventoryService) InventoryByMarketID(marketID string) (*core.Inventory, error) {
	return s.inventoryStg.GetByMarketID(marketID)
}

func (s *InventoryService) Set(_ context.Context, inv *core.Inventory) error {
	if err := inv.CheckCreate(); err != nil {
		return errors.New(core.InventoryErrRequiredFields, err)
	}

	defer func() {
		if _, err := s.marketStg.Index(inv.MarketID); err != nil {
			log.Printf("could not index market %s: %s", inv.MarketID, err)
		}
	}()

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
		return s.inventoryStg.Update(inv)
	}

	return s.inventoryStg.Create(inv)
}
