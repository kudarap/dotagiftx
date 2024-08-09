package steaminvorg

import (
	"reflect"
	"testing"
	"time"
)

func TestMetadata_isCacheFresh(t *testing.T) {
	tests := []struct {
		lastUpdated time.Time
		isFresh     bool
	}{
		{time.Time{}, false},
		{time.Date(2021, time.November, 10, 23, 0, 0, 0, time.UTC), false},
		{time.Date(2029, time.November, 10, 23, 0, 0, 0, time.UTC), true},
	}
	for _, tt := range tests {
		m := Metadata{LastUpdated: tt.lastUpdated}
		if got := m.isCacheFresh(); got != tt.isFresh {
			t.Errorf("isCacheFresh() = %v, want %v", got, tt.isFresh)
		}
	}
}

/*
{
  "Count": 1,
  "Items": [
    {
      "all_inventory_length": {
        "N": "0"
      },
      "all_description_length": {
        "N": "0"
      },
      "queued_timestamp": {
        "N": "1684026992770"
      },
      "profile": {
        "S": "76561198422287664"
      },
      "status": {
        "S": "error:403"
      },
      "timestamp": {
        "N": "1684026993336"
      },
      "result_url": {
        "S": "steam fetch error: profile is private"
      },
      "id": {
        "NULL": true
      },
      "inventory_format": {
        "N": "2"
      }
    }
  ],
  "ScannedCount": 1
}
*/

func Test_parseStrTime(t *testing.T) {
	tests := []struct {
		ts   string
		want time.Time
	}{
		{"", time.Time{}},
		{"invalid", time.Time{}},
		{"1684026992770", time.Date(2023, time.May, 14, 9, 16, 32, 0, time.Local)},
	}
	for _, tt := range tests {
		if got := parseStrTime(tt.ts); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("parseStrTime() = %v, want %v", got, tt.want)
		}
	}
}
