package redis

import (
	"testing"
)

func TestNew(t *testing.T) {
	c, err := New("localhost:6379", "root", 0)
	if c == nil {
		t.Errorf("New() got = %v, want not nil", c)
	}
	if err != nil {
		t.Errorf("New() error = %v, wantErr %v", err, nil)
	}
}
