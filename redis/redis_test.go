package redis

import (
	"testing"
)

func TestNew(t *testing.T) {
	c, err := New(Config{"localhost:6379", 0, "root"})
	if c == nil {
		t.Errorf("New() got = %v, want not nil", c)
	}
	if err != nil {
		t.Errorf("New() error = %v, wantErr %v", err, nil)
	}
}
