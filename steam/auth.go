package steam

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var (
	steamLogin = "https://steamcommunity.com/openid/login"

	openidMode       = "checkid_setup"
	openidNs         = "http://specs.openid.net/auth/2.0"
	openidIdentifier = "http://specs.openid.net/auth/2.0/identifier_select"

	validationRegexp       = regexp.MustCompile("^(http|https)://steamcommunity.com/openid/id/[0-9]{15,25}$")
	digitsExtractionRegexp = regexp.MustCompile("\\D+")
)

type OpenId struct {
	root      string
	returnUrl string
	data      url.Values
}

func NewOpenId(r *http.Request, cfg Config) *OpenId {
	id := new(OpenId)

	proto := "http://"
	if r.TLS != nil {
		proto = "https://"
	}
	id.root = proto + r.Host
	if cfg.Realm != "" {
		id.root = cfg.Realm
	}

	uri := r.RequestURI
	if i := strings.Index(uri, "openid"); i != -1 {
		uri = uri[0 : i-1]
	}
	id.returnUrl = id.root + uri
	if cfg.Return != "" {
		id.returnUrl = cfg.Return
	}

	switch r.Method {
	case "POST":
		id.data = r.Form
	case "GET":
		id.data = r.URL.Query()
	}

	return id
}

func (id OpenId) AuthUrl() string {
	data := map[string]string{
		"openid.claimed_id": openidIdentifier,
		"openid.identity":   openidIdentifier,
		"openid.mode":       openidMode,
		"openid.ns":         openidNs,
		"openid.realm":      id.root,
		"openid.return_to":  id.returnUrl,
	}

	i := 0
	u := steamLogin + "?"
	for key, value := range data {
		u += key + "=" + value
		if i != len(data)-1 {
			u += "&"
		}
		i++
	}
	return u
}

func (id *OpenId) ValidateAndGetId() (string, error) {
	if id.Mode() != "id_res" {
		return "", errors.New("Mode must equal to \"id_res\".")
	}

	if id.data.Get("openid.return_to") != id.returnUrl {
		return "", errors.New("The \"return_to url\" must match the url of current request.")
	}

	params := make(url.Values)
	params.Set("openid.assoc_handle", id.data.Get("openid.assoc_handle"))
	params.Set("openid.signed", id.data.Get("openid.signed"))
	params.Set("openid.sig", id.data.Get("openid.sig"))
	params.Set("openid.ns", id.data.Get("openid.ns"))

	split := strings.Split(id.data.Get("openid.signed"), ",")
	for _, item := range split {
		params.Set("openid."+item, id.data.Get("openid."+item))
	}
	params.Set("openid.mode", "check_authentication")

	resp, err := http.PostForm(steamLogin, params)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	response := strings.Split(string(content), "\n")
	if response[0] != "ns:"+openidNs {
		return "", errors.New("Wrong ns in the response.")
	}
	if strings.HasSuffix(response[1], "false") {
		return "", errors.New("Unable validate openId.")
	}

	openIdUrl := id.data.Get("openid.claimed_id")
	if !validationRegexp.MatchString(openIdUrl) {
		return "", errors.New("Invalid steam id pattern.")
	}

	return digitsExtractionRegexp.ReplaceAllString(openIdUrl, ""), nil
}

func (id OpenId) ValidateAndGetUser(apiKey string) (*PlayerSummaries, error) {
	steamId, err := id.ValidateAndGetId()
	if err != nil {
		return nil, err
	}
	return GetPlayerSummaries(steamId, apiKey)
}

func (id OpenId) Mode() string {
	return id.data.Get("openid.mode")
}
