package service

import (
	"context"
	"log"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/errors"
)

// NewInventory returns new Inventory service.
func NewInventory(rs dotagiftx.InventoryStorage, ms dotagiftx.MarketStorage, cs dotagiftx.CatalogStorage) dotagiftx.InventoryService {
	return &InventoryService{rs, ms, cs}
}

type InventoryService struct {
	inventoryStg dotagiftx.InventoryStorage
	marketStg    dotagiftx.MarketStorage
	catalogStg   dotagiftx.CatalogStorage
}

func (s *InventoryService) Inventories(opts dotagiftx.FindOpts) ([]dotagiftx.Inventory, *dotagiftx.FindMetadata, error) {
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

	return res, &dotagiftx.FindMetadata{
		ResultCount: len(res),
		TotalCount:  tc,
	}, nil
}

func (s *InventoryService) Inventory(id string) (*dotagiftx.Inventory, error) {
	inv, err := s.inventoryStg.Get(id)
	if err != nil && err != dotagiftx.InventoryErrNotFound {
		return nil, err
	}
	if inv != nil {
		return inv, nil
	}

	// If we can't find using id. lets try market ID
	return s.inventoryStg.GetByMarketID(id)
}

func (s *InventoryService) InventoryByMarketID(marketID string) (*dotagiftx.Inventory, error) {
	return s.inventoryStg.GetByMarketID(marketID)
}

func (s *InventoryService) Set(_ context.Context, inv *dotagiftx.Inventory) error {
	if err := inv.CheckCreate(); err != nil {
		return errors.New(dotagiftx.InventoryErrRequiredFields, err)
	}

	defer func() {
		mkt, err := s.marketStg.Index(inv.MarketID)
		if err != nil {
			log.Printf("could not index market %s: %s", inv.MarketID, err)
		}
		if _, err = s.catalogStg.Index(mkt.ItemID); err != nil {
			log.Printf("could not index catalog %s: %s", inv.MarketID, err)
		}
	}()

	// Update market Inventory status.
	if err := s.marketStg.BaseUpdate(&dotagiftx.Market{
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
