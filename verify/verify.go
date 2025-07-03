package verify

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/kudarap/dotagiftx/steam"
)

// AssetSource represents inventory asset source provider.
type AssetSource func(ctx context.Context, steamID string) ([]steam.Asset, error)

func MergeAssetSource(providers ...AssetSource) AssetSource {
	logger := slog.Default()
	return func(ctx context.Context, steamID string) ([]steam.Asset, error) {
		for name, source := range providers {
			logger.Info("attempting to fetch asset",
				"provider", name,
				"steam_id", steamID,
			)
			assets, err := source(ctx, steamID)
			if err != nil {
				logger.Error("failed to fetch assets",
					"provider", name,
					"steam_id", steamID,
					"err", err,
				)

				// stop the retrying when error is private.
				if errors.Is(err, steam.ErrInventoryPrivate) {
					logger.Info("stopping due to private inventory")
					return nil, err
				}

				continue
			}
			return assets, nil
		}
		return nil, fmt.Errorf("all asset providers attempted: %s", steamID)
	}
}
