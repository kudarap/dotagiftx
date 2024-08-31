package rethink

import "testing"

func Test_spanCleanUUIDs(t *testing.T) {
	tests := []struct {
		in, out string
	}{
		{"19176a5b-8361-42a2-9f32-b368fb3b46ce", "<uuid>"},
		{
			`rethink list r.Table("inventory").GetAll("1e7be262-2fc8-4496-b470-4ed1c195ac43", index="market_id")`,
			`rethink list r.Table("inventory").GetAll("<uuid>", index="market_id")`,
		},
	}
	for _, tt := range tests {
		if got := spanCleanUUIDs(tt.in); got != tt.out {
			t.Errorf("spanCleanUUIDs() = %v, want %v", got, tt.out)
		}
	}
}
