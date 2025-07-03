package verify

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/steam"
)

type Service struct {
	assetSources []AssetSourceContext
	inventorySvc dotagiftx.InventoryService
	deliverySvc  dotagiftx.DeliveryService
}

func NewService(
	as []AssetSourceContext,
	is dotagiftx.InventoryService,
	ds dotagiftx.DeliveryService,
) *Service {
	return &Service{as, is, ds}
}

// Inventory checks item existence on inventory and returns an error when request has status error or response
// body malformed.
func (s *Service) Inventory(ctx context.Context, marketID, steamID, itemName string) error {
	if steamID == "" || itemName == "" {
		return fmt.Errorf("all params are required")
	}

	var status dotagiftx.InventoryStatus
	source := s.assetProvider()
	provider, assets, err := source(ctx, steamID)
	if err != nil {
		status = dotagiftx.InventoryStatusError
		if errors.Is(err, steam.ErrInventoryPrivate) {
			status = dotagiftx.InventoryStatusPrivate
		}
	}
	if status != 0 {
		assets = filterByName(assets, itemName)
		assets = filterByGiftable(assets)
		status = dotagiftx.InventoryStatusNoHit
		if len(assets) != 0 {
			status = dotagiftx.InventoryStatusVerified
		}
	}

	err = s.inventorySvc.Set(ctx, &dotagiftx.Inventory{
		MarketID:   marketID,
		Status:     status,
		Assets:     assets,
		VerifiedBy: provider,
	})
	if err != nil {
		return fmt.Errorf("failed saving inventory: %w", err)
	}
	return nil
}

func (s *Service) Delivery(ctx context.Context, steamID, itemName string) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) assetProvider() AssetSourceContext {
	logger := slog.Default()
	return func(ctx context.Context, steamID string) (string, []steam.Asset, error) {
		for _, source := range s.assetSources {
			name, assets, err := source(ctx, steamID)
			if err != nil {
				logger.Error("failed to fetch assets",
					"provider", name, "steam_id", steamID,
					"err", err,
				)
				continue
			}
			return name, assets, nil
		}

		return "", nil, fmt.Errorf("all asset providers attempted: %s", steamID)
	}
}

// AssetSource represents inventory asset source provider.
type AssetSource func(steamID string) ([]steam.Asset, error)

type AssetSourceContext func(ctx context.Context, steamID string) (string, []steam.Asset, error)

func MultiAssetSource(providers ...AssetSource) AssetSource {
	logger := slog.Default()
	return func(steamID string) ([]steam.Asset, error) {
		for name, source := range providers {
			logger.Info("attempting to fetch asset",
				"provider", name,
				"steam_id", steamID,
			)
			assets, err := source(steamID)
			if err != nil {
				logger.Error("failed to fetch assets",
					"provider", name, "steam_id", steamID,
					"err", err,
				)

				// stop the retrying when error is private.
				if errors.Is(err, steam.ErrInventoryPrivate) {
					return nil, err
				}

				continue
			}
			return assets, nil
		}
		return nil, fmt.Errorf("all asset providers attempted: %s", steamID)
	}
}

// filterByName filters item that matches the name or in the description that supports unbundled items.
func filterByName(a []steam.Asset, itemName string) []steam.Asset {
	var matches []steam.Asset
	for _, asset := range a {
		// Strip "bundle" string to cover items that unbundled:
		// - Dipper the Destroyer Bundle
		// - The Abscesserator Bundle
		itemName = strings.TrimSpace(strings.TrimSuffix(itemName, "Bundle"))
		if !strings.Contains(strings.Join(asset.Descriptions, "|"), itemName) &&
			!strings.Contains(asset.Name, itemName) {
			continue
		}

		// Excluded golden variant of the item.
		if asset.IsGoldenVariant(itemName) {
			continue
		}

		matches = append(matches, asset)
	}
	return matches
}

// filterByGiftable filters item that can be gifted and allowed to be listed.
func filterByGiftable(a []steam.Asset) []steam.Asset {
	var matches []steam.Asset
	for _, aa := range a {
		// Is the item unbundled but giftable?
		//
		// Is the item immortal and does not say its giftable?
		// This fixes the removed "Gift once" string on description
		// recently by Valve.
		if !aa.GiftOnce && (aa.IsCollectorsCache() || !aa.IsImmortal()) {
			continue
		}

		matches = append(matches, aa)
	}
	return matches
}
