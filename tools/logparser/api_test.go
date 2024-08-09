package main

import (
	"reflect"
	"testing"
	"time"
)

func Test_parseLog(t *testing.T) {
	tests := []struct {
		raw  string
		want ApiLog
	}{
		{
			`time="2023-10-07T02:09:37Z" level=info msg="[xx/yy-00] "GET http://api.dotagiftx.com/catalogs HTTP/1.1" from 103.105.213.190 - 200 0B in 728µs"`,
			ApiLog{
				"GET /catalogs",
				200,
				time.Date(2023, 10, 07, 2, 9, 37, 0, time.UTC),
				728 * time.Microsecond,
			},
		},
	}
	for _, tt := range tests {
		if got := parseLog(tt.raw); !reflect.DeepEqual(got[0], tt.want) {
			t.Errorf("parseLog() = %v, want %v", got[0], tt.want)
		}
	}
}
