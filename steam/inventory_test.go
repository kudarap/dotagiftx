package steam

import (
	"os"
	"reflect"
	"testing"
)

var testAssetData = map[string]RawInventoryAsset{
	"100000000": {
		ID:         "100000000",
		ClassID:    "3305750400",
		InstanceID: "3307872803",
	},
}
var descGothicWhisper = RawInventoryDesc{
	ClassID:    "3305750400",
	InstanceID: "3307872803",
	Name:       "Gothic Whisper",
	Image:      "TESTDATA_LARGE_IMAGE",
	Type:       "Mythical Bundle",
	Descriptions: RawInventoryItemDetails{
		{"Used By: Phantom Assassin"},
		{"The International 2019"},
		{"Gift From: gippeum"},
		{"Date Received: Aug 24, 2020 (23:15:11)"},
	},
}
var descEmptyDetails = RawInventoryDesc{
	ClassID:      "3305750400",
	InstanceID:   "3307872803",
	Name:         "Gothic Whisper",
	Image:        "TESTDATA_LARGE_IMAGE",
	Type:         "Mythical Bundle",
	Descriptions: nil,
}
var descGothicWhisperUnWrapped = RawInventoryDesc{
	ClassID:    "3305750400",
	InstanceID: "3307872803",
	Name:       "Gothic Whisper",
	Image:      "TESTDATA_LARGE_IMAGE",
	Type:       "Mythical Mysterious Item",
	Descriptions: RawInventoryItemDetails{
		{"Contains: Gothic Whisper"},
		{"Gift From: gippeum"},
	},
}

var flatGothicWhisper = Asset{
	AssetID:      "100000000",
	Name:         "Gothic Whisper",
	Image:        "TESTDATA_LARGE_IMAGE",
	Type:         "Mythical Bundle",
	Hero:         "Phantom Assassin",
	GiftFrom:     "gippeum",
	DateReceived: "Aug 24, 2020 (23:15:11)",
	Descriptions: []string{
		"Used By: Phantom Assassin",
		"The International 2019",
		"Gift From: gippeum",
		"Date Received: Aug 24, 2020 (23:15:11)",
	},
}

func Test_inventoryParser(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *RawInventory
		wantErr bool
	}{
		{
			"good base model",
			args{"./testdata/base_model.json"},
			&RawInventory{
				Success:   true,
				More:      false,
				MoreStart: 0,
				Assets:    testAssetData,
				Descriptions: map[string]RawInventoryDesc{
					"3305750400_3307872803": descGothicWhisper,
				},
			},
			false,
		},
		{
			"paginated",
			args{"./testdata/paginated.json"},
			&RawInventory{
				Success:   true,
				More:      false,
				MoreStart: 1,
				Assets:    testAssetData,
				Descriptions: map[string]RawInventoryDesc{
					"3305750400_3307872803": descGothicWhisper,
				},
			},
			false,
		},
		{
			"empty item description",
			args{"./testdata/empty_desc.json"},
			&RawInventory{
				Success:   true,
				More:      false,
				MoreStart: 0,
				Assets:    testAssetData,
				Descriptions: map[string]RawInventoryDesc{
					"3305750400_3307872803": descEmptyDetails,
				},
			},
			false,
		},
		{
			"unwrapped gift",
			args{"./testdata/unwrapped_gift.json"},
			&RawInventory{
				Success:   true,
				More:      false,
				MoreStart: 0,
				Assets:    testAssetData,
				Descriptions: map[string]RawInventoryDesc{
					"3305750400_3307872803": descGothicWhisperUnWrapped,
				},
			},
			false,
		},
		{
			"success false",
			args{"./testdata/success_false.json"},
			&RawInventory{
				Success: false,
			},
			false,
		},
		{
			"private inventory",
			args{"./testdata/private.json"},
			&RawInventory{
				Success: false,
				Error:   "This profile is private.",
			},
			false,
		},
		{
			"bad filepath",
			args{"./testdata/badfilepath.json"},
			nil,
			true,
		},
		{
			"bad json or malformed",
			args{"./testdata/malformed.json"},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := os.Open(tt.args.path)
			got, err := inventoryParser(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("newInventoryFromFile() \nerror = %v, \nwantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newInventoryFromFile() \n\ngot  %#v, \n\nwant %#v\n\n", got, tt.want)
			}
		})
	}
}

func Test_assetParser(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []Asset
		wantErr bool
	}{
		{
			"base model",
			args{"./testdata/base_model.json"},
			[]Asset{
				flatGothicWhisper,
			},
			false,
		},
		{
			"empty item description",
			args{"./testdata/empty_desc.json"},
			[]Asset{
				{
					AssetID: "100000000",
					Name:    "Gothic Whisper",
					Image:   "TESTDATA_LARGE_IMAGE",
					Type:    "Mythical Bundle",
				},
			},
			false,
		},
		{
			"bad filepath",
			args{"./testdata/badfilepath.json"},
			nil,
			true,
		},
		{
			"bad json or malformed",
			args{"./testdata/malformed.json"},
			nil,
			true,
		},
		// TODO: parse error
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := os.Open(tt.args.path)
			got, err := assetParser(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("newFlatInventoryFromFile() \nerror = %v, \nwantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newFlatInventoryFromFile() \n\ngot  %v, \n\nwant %v", got, tt.want)
			}
		})
	}
}
