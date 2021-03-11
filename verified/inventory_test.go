package verified

import (
	"testing"

	"github.com/kudarap/dotagiftx/steaminv"
)

func TestVerifyInventory(t *testing.T) {
	type args struct {
		steamID  string
		itemName string
	}
	tests := []struct {
		name    string
		args    args
		want    VerifyStatus
		count   int
		wantErr bool
	}{
		{"ok item", args{
			"76561198355627060",
			"Shattered Greatsword",
		}, VerifyStatusItem, 1, false},
		{"no hit", args{
			"76561198042690669",
			"Baby Demon",
		}, VerifyStatusNoHit, 0, false},
		{"private data", args{
			"76561198011477544",
			"Baby Demon",
		}, VerifyStatusPrivate, 0, false},
		{"bad steam id", args{
			"76561_____477544",
			"Bad Demon",
		}, VerifyStatusError, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, assets, err := Inventory(steaminv.InventoryAsset, tt.args.steamID, tt.args.itemName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Inventory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Inventory() got = %v, want %v", got, tt.want)
			}
			if len(assets) != tt.count {
				t.Errorf("Inventory() count = %v, want %v", len(assets), tt.count)
			}
		})
	}
}
