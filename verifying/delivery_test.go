package verifying

import (
	"reflect"
	"testing"

	dgx "github.com/kudarap/dotagiftx"
	"github.com/kudarap/dotagiftx/steam"
	"github.com/kudarap/dotagiftx/steaminvorg"
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
		want    dgx.DeliveryStatus
		count   int
		wantErr bool
	}{
		{"seller verified", args{
			"kudarap",
			"76561198287849998",
			"Sylvan Vedette",
		}, dgx.DeliveryStatusSenderVerified, 1, false},
		{"item verified", args{
			"Berserk",
			"76561198355627060",
			"Shattered Greatsword",
		}, dgx.DeliveryStatusNameVerified, 1, false},
		{"no hit", args{
			"kudarap",
			"76561198042690669",
			"Baby Demon",
		}, dgx.DeliveryStatusNoHit, 0, false},
		{"private data", args{
			"kudarap",
			"76561198011477544",
			"Baby Demon",
		}, dgx.DeliveryStatusPrivate, 0, false},
		{"bad steam id", args{
			"kudarap",
			"76561198011477544_",
			"Bad Demon",
		}, dgx.DeliveryStatusError, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, assets, err := Delivery(steaminvorg.InventoryAsset, tt.args.sellerPersona, tt.args.buyerSteamID, tt.args.itemName)
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
			"76561198287849998",
			"Sylvan Vedette",
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stat1, assets1, err1 := Delivery(steaminvorg.InventoryAsset, tt.args.sellerPersona, tt.args.buyerSteamID, tt.args.itemName)
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
