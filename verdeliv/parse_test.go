package verdeliv

import (
	"reflect"
	"testing"
)

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
		// TODO: Add test cases.
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
					"3305750400_3307872803": {
						ClassID:    "3305750400",
						InstanceID: "3307872803",
						Name:       "Gothic Whisper",
						Image:      "TESTDATA_LARGE_IMAGE",
						Type:       "Mythical Bundle",
						Descriptions: []itemDetails{
							{"Used By: Phantom Assassin"},
							{"The International 2019"},
						},
					},
				},
			},
			false,
		},
		// TODO: valid empty inventory
		// TODO: sample inventory
		// TODO: private inventory
		// TODO: bad inventory
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
