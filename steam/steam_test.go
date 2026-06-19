package steam

import "testing"

func Test_cleanProfileURL(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{
			"https://steamcommunity.com/profiles/76561198088587178",
			"https://steamcommunity.com/profiles/76561198088587178",
		},
		{
			"https://steamcommunity.com/profiles/76561198088587178/",
			"https://steamcommunity.com/profiles/76561198088587178",
		},
		{
			"https://steamcommunity.com/profiles/76561198088587178/inventory",
			"https://steamcommunity.com/profiles/76561198088587178",
		},
		{
			"https://steamcommunity.com/profiles/76561198088587178/inventory/#570",
			"https://steamcommunity.com/profiles/76561198088587178",
		},
	}
	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			if got := cleanProfileURL(tt.url); got != tt.want {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}
