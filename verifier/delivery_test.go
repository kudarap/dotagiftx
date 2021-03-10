package verifier

import (
	"testing"
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
		// TODO: Add test cases.
		{"ok kudarap", args{
			"kudarap",
			"76561198042690669",
			"Riddle of the Hierophant",
		}, VerifyStatusSeller, 5, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, assets, err := VerifyDelivery(tt.args.sellerPersona, tt.args.buyerSteamID, tt.args.itemName)
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
