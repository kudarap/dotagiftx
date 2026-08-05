package steamcrawl

import (
	"strings"
	"testing"
)

func TestNewSession_headers(t *testing.T) {
	h := NewSession("76561198088587178")

	for _, key := range []string{
		"User-Agent", "Accept", "Accept-Language", "Accept-Encoding",
		"Connection", "Referer", "Sec-Fetch-Dest", "Sec-Fetch-Mode",
		"Sec-Fetch-Site", "Sec-GPC", "Pragma", "Cache-Control", "Cookie",
	} {
		if h.Get(key) == "" {
			t.Fatalf("missing header %s", key)
		}
	}

	if !strings.Contains(h.Get("Referer"), "76561198088587178") {
		t.Fatalf("referer = %q, want steam id", h.Get("Referer"))
	}
	if !strings.Contains(h.Get("Cookie"), "sessionid=") {
		t.Fatalf("cookie = %q, want sessionid", h.Get("Cookie"))
	}
}

func TestNewSession_freshEachCall(t *testing.T) {
	a := NewSession("1")
	b := NewSession("1")
	if a.Get("Cookie") == b.Get("Cookie") {
		t.Fatalf("cookie not randomized: %q", a.Get("Cookie"))
	}
}

func TestNewSession_variedUserAgent(t *testing.T) {
	seen := map[string]bool{}
	for i := 0; i < 100; i++ {
		seen[NewSession("1").Get("User-Agent")] = true
		if len(seen) > 1 {
			return
		}
	}
	t.Fatalf("user agent never varied across 100 calls: %v", seen)
}

func Test_freshSessionCookie(t *testing.T) {
	c := freshSessionCookie()
	if !strings.Contains(c, "sessionid=") {
		t.Fatalf("cookie = %q, want sessionid", c)
	}
}
