package service

import (
	"context"

	"github.com/kudarap/dotagiftx/core"
	"github.com/kudarap/dotagiftx/errors"
)

// NewDelivery returns new Delivery service.
func NewDelivery(rs core.DeliveryStorage, ms core.MarketStorage) core.DeliveryService {
	return &deliveryService{rs, ms}
}

type deliveryService struct {
	deliveryStg core.DeliveryStorage
	marketStg   core.MarketStorage
}

func (s *deliveryService) Deliveries(opts core.FindOpts) ([]core.Delivery, *core.FindMetadata, error) {
	res, err := s.deliveryStg.Find(opts)
	if err != nil {
		return nil, nil, err
	}

	if !opts.WithMeta {
		return res, nil, err
	}

	// Get result and total count for metadata.
	tc, err := s.deliveryStg.Count(opts)
	if err != nil {
		return nil, nil, err
	}

	return res, &core.FindMetadata{
		ResultCount: len(res),
		TotalCount:  tc,
	}, nil
}

func (s *deliveryService) Delivery(id string) (*core.Delivery, error) {
	inv, err := s.deliveryStg.Get(id)
	if err != nil && err != core.DeliveryErrNotFound {
		return nil, err
	}
	if inv != nil {
		return inv, nil
	}

	// If we can't find using id. lets try market ID
	return s.deliveryStg.GetByMarketID(id)
}

func (s *deliveryService) DeliveryByMarketID(marketID string) (*core.Delivery, error) {
	return s.deliveryStg.GetByMarketID(marketID)
}

func (s *deliveryService) Set(_ context.Context, del *core.Delivery) error {
	if err := del.CheckCreate(); err != nil {
		return errors.New(core.DeliveryErrRequiredFields, err)
	}

	// Detect if there are still un-opened gift.
	del = del.IsGiftOpened()

	// Update market delivery status.
	if err := s.marketStg.BaseUpdate(&core.Market{
		ID:             del.MarketID,
		DeliveryStatus: del.Status,
	}); err != nil {
		return err
	}

	// Just update existing record.
	cur, _ := s.DeliveryByMarketID(del.MarketID)
	if cur != nil {
		del.ID = cur.ID
		del.Retries = cur.Retries + 1
		del = del.AddAssets(cur.Assets)
		return s.deliveryStg.Update(del)
	}

	return s.deliveryStg.Create(del)
}
