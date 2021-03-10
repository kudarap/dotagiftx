package verifier

import (
	"testing"

	"github.com/kudarap/dotagiftx/steaminv"
)

func TestVerifyDelivery(t *testing.T) {
	type args struct {
		sellerPersona string
		buyerSteamID  string
		itemName      string
	}
	tests := []struct {
		name    string
		args    args
		want    VerifyStatus
		count   int
		wantErr bool
	}{
		{"ok seller", args{
			"kudarap",
			"76561198042690669",
			"Riddle of the Hierophant",
		}, VerifyStatusSeller, 5, false},
		{"ok item", args{
			"Berserk",
			"76561198355627060",
			"Shattered Greatsword",
		}, VerifyStatusItem, 1, false},
		{"no hit", args{
			"kudarap",
			"76561198042690669",
			"Baby Demon",
		}, VerifyStatusNoHit, 0, false},
		{"private data", args{
			"kudarap",
			"76561198011477544",
			"Baby Demon",
		}, VerifyStatusPrivate, 0, false},
		{"bad steam id", args{
			"kudarap",
			"76561198011477544_",
			"Bad Demon",
		}, VerifyStatusError, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, assets, err := Delivery(steaminv.InventoryAsset, tt.args.sellerPersona, tt.args.buyerSteamID, tt.args.itemName)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyDelivery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("VerifyDelivery() got = %v, want %v", got, tt.want)
			}
			if len(assets) != tt.count {
				t.Errorf("VerifyDelivery() count = %v, want %v", len(assets), tt.count)
			}
		})
	}
}
