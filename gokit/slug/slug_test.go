package slug

import "testing"

func TestMake(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"", ""},
		{"Hello World", "hello-world"},
		{"chicken & beer", "chicken-beer"},
		{"float 64", "float-64"},
		{"go kit slug", "gokit-slug"},
		{"gokit-slug", "gokit-slug"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := Make(tt.in); got != tt.out {
				t.Errorf("Make() = %v, want %v", got, tt.out)
			}
		})
	}
}
