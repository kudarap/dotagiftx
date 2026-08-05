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
		{
			"https://steamcommunity.com/id/kudrap",
			"https://steamcommunity.com/id/kudrap",
		},
		{
			"https://steamcommunity.com/id/kudrap/inventory",
			"https://steamcommunity.com/id/kudrap",
		},
		{
			"https://steamcommunity.com/id/kudrap/inventory/#570",
			"https://steamcommunity.com/id/kudrap",
		},
	}
	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			if got, _ := cleanProfileURL(tt.url); got != tt.want {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func TestValidateSteamID(t *testing.T) {
	tests := []struct {
		steamID string
		wantErr bool
	}{
		{"", true},
		{"http:", true},
		{"76561198068062691/inventory", true},
		{"76561198068062691", false},
	}
	for _, tt := range tests {
		t.Run(tt.steamID, func(t *testing.T) {
			if err := ValidateSteamID(tt.steamID); (err != nil) != tt.wantErr {
				t.Errorf("ValidateSteamID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
