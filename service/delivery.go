package service

import (
	"context"
	"log"

	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/errors"
)

// NewDelivery returns new Delivery service.
func NewDelivery(rs dgx.DeliveryStorage, ms dgx.MarketStorage) dgx.DeliveryService {
	return &deliveryService{rs, ms}
}

type deliveryService struct {
	deliveryStg dgx.DeliveryStorage
	marketStg   dgx.MarketStorage
}

func (s *deliveryService) Deliveries(opts dgx.FindOpts) ([]dgx.Delivery, *dgx.FindMetadata, error) {
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

	return res, &dgx.FindMetadata{
		ResultCount: len(res),
		TotalCount:  tc,
	}, nil
}

func (s *deliveryService) Delivery(id string) (*dgx.Delivery, error) {
	inv, err := s.deliveryStg.Get(id)
	if err != nil && err != dgx.DeliveryErrNotFound {
		return nil, err
	}
	if inv != nil {
		return inv, nil
	}

	// If we can't find using id. lets try market ID
	return s.deliveryStg.GetByMarketID(id)
}

func (s *deliveryService) DeliveryByMarketID(marketID string) (*dgx.Delivery, error) {
	return s.deliveryStg.GetByMarketID(marketID)
}

func (s *deliveryService) Set(_ context.Context, del *dgx.Delivery) error {
	if err := del.CheckCreate(); err != nil {
		return errors.New(dgx.DeliveryErrRequiredFields, err)
	}

	defer func() {
		if _, err := s.marketStg.Index(del.MarketID); err != nil {
			log.Printf("could not index market %s: %s", del.MarketID, err)
		}
	}()

	// Detect if there are still un-opened gift.
	del = del.IsGiftOpened()

	// Update market delivery status.
	if err := s.marketStg.BaseUpdate(&dgx.Market{
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
