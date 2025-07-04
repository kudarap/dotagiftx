package verify

import (
	"context"
	"errors"
	"fmt"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/steam"
)

// AssetSource represents inventory asset source provider.
type AssetSource func(ctx context.Context, steamID string) (providerID string, sa []steam.Asset, err error)

type Source struct {
	providers []AssetSource
}

func NewSource(as ...AssetSource) *Source {
	return &Source{as}
}

type InventorResult struct {
	Status     dotagiftx.InventoryStatus
	Assets     []steam.Asset
	VerifiedBy string
}

func (s *Source) Inventory(ctx context.Context, steamID, itemName string) (*InventorResult, error) {
	src := JoinAssetSource(s.providers...)
	res, err := Inventory(ctx, src, steamID, itemName)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type DeliveryResult struct {
	Status     dotagiftx.DeliveryStatus
	Assets     []steam.Asset
	VerifiedBy string
}

func (s *Source) Delivery(ctx context.Context, sellerPersona, steamID, itemName string) (*DeliveryResult, error) {
	src := JoinAssetSource(s.providers...)
	res, err := Delivery(ctx, src, sellerPersona, steamID, itemName)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func JoinAssetSource(providers ...AssetSource) AssetSource {
	return func(ctx context.Context, steamID string) (string, []steam.Asset, error) {
		for _, source := range providers {
			name, assets, err := source(ctx, steamID)
			if err != nil {
				if errors.Is(err, steam.ErrInventoryPrivate) {
					return name, nil, err
				}
				continue
			}
			return name, assets, nil
		}
		return "", nil, fmt.Errorf("all source exhausted: %s", steamID)
	}
}
