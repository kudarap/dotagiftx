package dotagiftx

import "testing"

func Test_makeSlug(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"", ""},
		{"Hello World", "hello-world"},
		{"chicken & beer", "chicken-beer"},
		{"float 64", "float-64"},
		{"go kit slug", "go-kit-slug"},
		{"gokit-slug", "gokit-slug"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := makeSlug(tt.in); got != tt.out {
				t.Errorf("makeSlug() = %v, want %v", got, tt.out)
			}
		})
	}
}
