package verify

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/kudarap/dotagiftx"
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
		want    dotagiftx.DeliveryStatus
		count   int
		wantErr bool
	}{
		{"seller verified", args{
			"kudarap",
			"76561198287849998",
			"Sylvan Vedette",
		}, dotagiftx.DeliveryStatusSenderVerified, 1, false},
		{"item verified", args{
			"Berserk",
			"76561198355627060",
			"Shattered Greatsword",
		}, dotagiftx.DeliveryStatusNameVerified, 1, false},
		{"no hit", args{
			"kudarap",
			"76561198042690669",
			"Baby Demon",
		}, dotagiftx.DeliveryStatusNoHit, 0, false},
		{"private data", args{
			"kudarap",
			"76561198011477544",
			"Baby Demon",
		}, dotagiftx.DeliveryStatusPrivate, 0, false},
		{"bad steam id", args{
			"kudarap",
			"76561198011477544_",
			"Bad Demon",
		}, dotagiftx.DeliveryStatusError, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip()

			ctx := context.Background()
			res, err := Delivery(
				ctx,
				steaminvorg.InventoryAssetWithProvider,
				tt.args.sellerPersona,
				tt.args.buyerSteamID,
				tt.args.itemName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delivery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got := res.Status
			assets := res.Assets
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
			t.Skip()

			ctx := context.Background()
			res1, err1 := Delivery(
				ctx,
				steaminvorg.InventoryAssetWithProvider,
				tt.args.sellerPersona,
				tt.args.buyerSteamID,
				tt.args.itemName,
			)
			res2, err2 := Delivery(
				ctx,
				steam.InventoryAssetWithProvider,
				tt.args.sellerPersona,
				tt.args.buyerSteamID,
				tt.args.itemName,
			)

			if !errors.Is(err2, err1) {
				t.Errorf("Delivery() error not matched %v x %v", err1, err2)
			}
			if res1.Status != res2.Status {
				t.Errorf("Delivery() status not matched %v x %v", res1.Status, res2.Status)
			}
			if !reflect.DeepEqual(res1.Assets, res2.Assets) {
				t.Errorf("Delivery() assets not matched \n\n%#v \n\n%#v", res1.Assets, res2.Assets)
			}
		})
	}
}
