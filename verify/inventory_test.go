package verify

import (
	"context"
	"testing"

	"github.com/kudarap/dotagiftx"
)

func TestVerifyInventory(t *testing.T) {
	tests := []struct {
		name     string
		steamID  string
		itemName string
		want     dotagiftx.InventoryStatus
		count    int
		wantErr  bool
	}{
		{
			"item verified",
			"76561198088587178",
			"Echoes of the Everblack",
			dotagiftx.InventoryStatusVerified,
			1,
			false,
		},
		{
			"no hit",
			"76561198042690669",
			"Baby Demon",
			dotagiftx.InventoryStatusNoHit,
			0,
			false,
		},
		{
			"private data",
			"76561198011477544",
			"Baby Demon",
			dotagiftx.InventoryStatusPrivate,
			0,
			false,
		},
		{
			"bad steam id",
			"76561_____477544",
			"Bad Demon",
			dotagiftx.InventoryStatusError,
			0,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip()

			ctx := context.Background()
			src := JoinAssetSource()
			res, err := Inventory(ctx, src, tt.steamID, tt.itemName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Inventory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got := res.Status
			assets := res.Assets
			if got != tt.want {
				t.Errorf("Inventory() got = %v, want %v", got, tt.want)
			}
			if len(assets) != tt.count {
				t.Errorf("Inventory() count = %v, want %v", len(assets), tt.count)
			}
		})
	}
}

func Test_fixMisspelledName(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want string
	}{
		{
			"backward compatible",
			"Intergalactic Orbliterator",
			"Intergalactic Orbliterator",
			"Intergalactic Orbliterator",
		},
		{
			"fix",
			"Intergalactic Orbliterator",
			"Intergalactic Obliterator",
			"Intergalactic Obliterator",
		},
		{
			"fix parts",
			"Intergalactic Orbliterator - Head",
			"Intergalactic Obliterator - Head",
			"Intergalactic Obliterator - Head",
		},
		{
			"skip fix",
			"Chrononaut Continuum",
			"Intergalactic Obliterator",
			"Chrononaut Continuum",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fixMisspelledName(tt.a, tt.b); got != tt.want {
				t.Errorf("fixMisspelledName() = %v, want %v", got, tt.want)
			}
		})
	}
}
