package phantasm

import (
	"reflect"
	"testing"
)

func Test_merge_verify(t *testing.T) {
	type args struct {
		res []*inventory
	}
	tests := []struct {
		name string
		args args
		want *inventory
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := merge(tt.args.res); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("merge() = %v, want %v", got, tt.want)
			}
		})
	}
}
