package steaminvorg

import (
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
