package verify

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/steam"
)

// Inventory checks item existence on inventory.
//
// Returns an error when request has status error or response body malformed.
func Inventory(ctx context.Context, source AssetSource, steamID, itemName string) (*InventorResult, error) {
	result := InventorResult{
		Status: dotagiftx.InventoryStatusError,
	}
	if steamID == "" || itemName == "" {
		return &result, fmt.Errorf("all params are required")
	}

	verifier, assets, err := source(ctx, steamID)
	result.VerifiedBy = verifier
	if err != nil {
		if errors.Is(err, steam.ErrInventoryPrivate) {
			result.Status = dotagiftx.InventoryStatusPrivate
			return &result, nil
		}
		return nil, err
	}

	assets = filterByName(assets, itemName)
	assets = filterByGiftable(assets)
	if len(assets) == 0 {
		result.Status = dotagiftx.InventoryStatusNoHit
		return &result, nil
	}

	result.Assets = assets
	result.Status = dotagiftx.InventoryStatusVerified
	return &result, nil
}

// filterByName filters item that matches the name or in the description that supports unbundled items.
func filterByName(a []steam.Asset, name string) []steam.Asset {
	var matches []steam.Asset
	for _, asset := range a {
		// Strip "bundle" suffix to cover unbundled items:
		// - Dipper the Destroyer Bundle
		// - The Abscesserator Bundle
		name = strings.TrimSpace(strings.TrimSuffix(name, "Bundle"))
		// Fix misspelled names like earth shakers arcana.
		asset.Name = fixMisspelledName(asset.Name, name)
		desc := fixMisspelledName(strings.Join(asset.Descriptions, "|"), name)
		if !strings.Contains(desc, name) && !strings.Contains(asset.Name, name) {
			continue
		}

		// Excluded golden variant of the item.
		if asset.IsGoldenVariant(name) {
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

func fixMisspelledName(a, b string) string {
	if strings.EqualFold(a, b) {
		return a
	}
	if strings.EqualFold(a, "Intergalactic Orbliterator") {
		return "Intergalactic Obliterator"
	}
	if strings.Contains(a, "Orbliterator") && strings.Contains(b, "Obliterator") {
		return strings.ReplaceAll(a, "Orbliterator", "Obliterator")
	}
	return a
}
