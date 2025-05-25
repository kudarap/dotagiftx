package dgx

import (
	"testing"
	"time"
)

func TestChat_Signature(t *testing.T) {
	tests := []struct {
		chat    Chat
		secret  []byte
		want    string
		wantErr bool
	}{
		{
			Chat{
				ID: "123",
				Users: []User{
					{ID: "user-1", SteamID: "7000", Name: "kudarap"},
					{ID: "user-2", SteamID: "7001", Name: "middaydemon"},
				},
				Item:     Item{ID: "item-1", Name: "Cannonroar Confessor"},
				MarketID: "market-1",
				Messages: []Message{
					{
						UserID:    "user-31",
						Content:   "Yo?",
						Timestamp: time.Date(2025, 5, 25, 7, 0, 0, 0, time.UTC),
					},
					{
						UserID:    "user-2",
						Content:   "Let's sprint forward!",
						Timestamp: time.Date(2025, 5, 25, 7, 1, 0, 0, time.UTC),
					},
				},
			},
			[]byte(`spring-forward`),
			"8b9a898f7d8409574b2f8a592dfcad84a67add4c8d645f7d3cb992df1756f7c7",
			false,
		},
		{
			Chat{
				ID: "123",
				Users: []User{{
					ID:      "user-1",
					SteamID: "7000",
					Name:    "kudarap",
				}},
				Item:     Item{ID: "item-1", Slug: "item/1"},
				MarketID: "market-1",
				Messages: []Message{
					{
						UserID:    "user-2",
						SteamID:   "7000",
						Content:   "yo",
						Timestamp: time.Time{},
					},
				},
			},
			[]byte(`spring-forward`),
			"892e2acdf37d20b35f562bf0dc6b812be37c0df1949968079b5b494811ccedc0",
			false,
		},
		{
			Chat{},
			[]byte(`secret`),
			"920eaa595171cf347bdbcff5c38410cb4019adfc3d618918b1fc970f58691a55",
			false,
		},
	}
	for _, tt := range tests {
		got, err := tt.chat.Signature(tt.secret)
		if (err != nil) != tt.wantErr {
			t.Errorf("Signature() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if got != tt.want {
			t.Errorf("Signature() got = %v, want %v", got, tt.want)
		}
	}
}

func TestChatVerify(t *testing.T) {
	tests := []struct {
		// params
		secret string
		file   string
		// returns
		want bool
	}{
		{
			`secret`,
			`sha256: 8b9a898f7d8409574b2f8a592dfcad84a67add4c8d645f7d3cb992df1756f7c7

market-ref: market-1

76561198088587178(kudarap): 25 May 25 07:00 UTC  
Yo?

76561198088587000(middaydemon): 25 May 25 07:01 UTC 
Let's sprint forward!
`,
			true,
		},
	}
	for _, tt := range tests {
		if got := ChatVerify([]byte(tt.secret), []byte(tt.file)); got != tt.want {
			t.Errorf("ChatVerify() = %v, want %v", got, tt.want)
		}
	}
}
