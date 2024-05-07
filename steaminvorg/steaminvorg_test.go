package steaminvorg

import (
	"reflect"
	"testing"
	"time"
)

func TestMetadata_isCacheFresh(t *testing.T) {
	tests := []struct {
		lastUpdated *time.Time
		isFresh     bool
	}{
		{nil, false},
	}
	for _, tt := range tests {
		m := Metadata{LastUpdated: tt.lastUpdated}
		if got := m.isCacheFresh(); got != tt.isFresh {
			t.Errorf("isCacheFresh() = %v, want %v", got, tt.isFresh)
		}
	}
}

func Test_parseStrTime(t *testing.T) {
	tests := []struct {
		ts   string
		want *time.Time
	}{
		{"", nil},
	}
	for _, tt := range tests {
		if got := parseStrTime(tt.ts); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("parseStrTime() = %v, want %v", got, tt.want)
		}
	}
}
