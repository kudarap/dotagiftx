package verifier

import (
	"reflect"
	"testing"

	"github.com/kudarap/dotagiftx/steam"
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
				t.Errorf("Delivery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Delivery() got = %v, want %v", got, tt.want)
			}
			if len(assets) != tt.count {
				t.Errorf("Delivery() count = %v, want %v", len(assets), tt.count)
			}
		})
	}
}

func TestVerifyDeliveryMultiSources(t *testing.T) {
	type args struct {
		sellerPersona string
		buyerSteamID  string
		itemName      string
	}
	tests := []struct {
		name string
		args args
	}{
		{"ok seller", args{
			"kudarap",
			"76561198073410102",
			"Aspect of Oscilla",
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stat1, assets1, err1 := Delivery(steaminv.InventoryAsset, tt.args.sellerPersona, tt.args.buyerSteamID, tt.args.itemName)
			stat2, assets2, err2 := Delivery(steam.InventoryAsset, tt.args.sellerPersona, tt.args.buyerSteamID, tt.args.itemName)

			if err1 != err2 {
				t.Errorf("Delivery() error not matched %v x %v", err1, err2)
			}

			if stat1 != stat2 {
				t.Errorf("Delivery() status not matched %v x %v", stat1, stat2)
			}

			if !reflect.DeepEqual(assets1, assets2) {
				t.Errorf("Delivery() assets not matched \n\n%#v \n\n%#v", assets1, assets2)
			}
		})
	}
}
