package steam

import (
	"fmt"
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type PlayerSummaries struct {
	SteamId                  string `json:"steamid"`
	CommunityVisibilityState int    `json:"communityvisibilitystate"`
	ProfileState             int    `json:"profilestate"`
	PersonaName              string `json:"personaname"`
	LastLogOff               int    `json:"lastlogoff"`
	ProfileUrl               string `json:"profileurl"`
	Avatar                   string `json:"avatar"`
	AvatarMedium             string `json:"avatarmedium"`
	AvatarFull               string `json:"avatarfull"`
	PersonaState             int    `json:"personastate"`

	CommentPermission int    `json:"commentpermission"`
	RealName          string `json:"realname"`
	//PrimaryClanId     string `json:"primaryclanid"`
	//TimeCreated       int    `json:"timecreated"`
	//LocCountryCode    string `json:"loccountrycode"`
	//LocStateCode      string `json:"locstatecode"`
	//LocCityId         int    `json:"loccityid"`
	//GameId            string `json:"gameid"`
	//GameExtraInfo     string `json:"gameextrainfo"`
	//GameServerIp      string `json:"gameserverip"`
}

func GetPlayerSummaries(steamId, apiKey string) (*PlayerSummaries, error) {
	url := fmt.Sprintf("https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s", apiKey, steamId)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type Result struct {
		Response struct {
			Players []PlayerSummaries `json:"players"`
		} `json:"response"`
	}
	var data Result
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	if len(data.Response.Players) == 0 {
		return nil, fmt.Errorf("no result")
	}

	return &data.Response.Players[0], err
}

func ResolveVanityURL(vanityURL, apiKey string) (steamID string, err error) {
	url := fmt.Sprintf("https://api.steampowered.com/ISteamUser/ResolveVanityURL/v1/?key=%s&vanityurl=%s", apiKey, vanityURL)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	type Result struct {
		Response struct {
			SteamID string `json:"steamid"`
			Success int    `json:"success"`
		} `json:"response"`
	}
	var data Result
	if err = json.Unmarshal(body, &data); err != nil {
		return
	}

	if data.Response.Success != 1 {
		err = fmt.Errorf("steam resolve vanity url error %d", data.Response.Success)
		return
	}

	return data.Response.SteamID, nil
}
