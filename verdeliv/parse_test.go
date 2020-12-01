package verdeliv

import (
	"reflect"
	"testing"
)

var descGothicWhisper = description{
	ClassID:    "3305750400",
	InstanceID: "3307872803",
	Name:       "Gothic Whisper",
	Image:      "TESTDATA_LARGE_IMAGE",
	Type:       "Mythical Bundle",
	Descriptions: []itemDetails{
		{"Used By: Phantom Assassin"},
		{"The International 2019"},
		{"Gift From: gippeum"},
		{"Date Received: Aug 24, 2020 (23:15:11)"},
	},
}

var flatGothicWhisper = flatInventory{
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

func Test_newInventoryFromFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *inventory
		wantErr bool
	}{
		{
			"testing base inventory",
			args{"./testdata/basemodel.json"},
			&inventory{
				Success:   true,
				More:      false,
				MoreStart: false,
				Assets: map[string]asset{
					"100000000": {
						ID:         "100000000",
						ClassID:    "3305750400",
						InstanceID: "3307872803",
					},
				},
				Descriptions: map[string]description{
					"3305750400_3307872803": descGothicWhisper,
				},
			},
			false,
		},
		// TODO: valid empty inventory
		// TODO: sample inventory
		// TODO: private inventory
		// TODO: bad filepath
		// TODO: bad json or malformed
		// TODO: success false
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newInventoryFromFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("newInventoryFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newInventoryFromFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newFlatInventoryFromFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []flatInventory
		wantErr bool
	}{
		// TODO: base model
		{
			"base model",
			args{"./testdata/basemodel.json"},
			[]flatInventory{
				flatGothicWhisper,
			},
			false,
		},
		// TODO: parse error
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newFlatInventoryFromFile(tt.args.path)
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
