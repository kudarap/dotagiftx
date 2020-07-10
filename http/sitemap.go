package http

import (
	"net/http"

	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"
	"github.com/kudarap/dota2giftables/core"
)

func handleSitemap(marketSvc core.MarketService, userSvc core.UserService) http.HandlerFunc {
	sm := stm.NewSitemap(1)
	//sm.SetVerbose(false)
	sm.SetDefaultHost("https://dota2giftables.com")
	sm.Create()

	// Static pages locations.
	sm.Add(stm.URL{{"loc", "/"}, {"changefreq", "daily"}, {"priority", 0.8}})
	sm.Add(stm.URL{{"loc", "/search"}, {"changefreq", "daily"}, {"priority", 0.6}})
	sm.Add(stm.URL{{"loc", "/search?sort=view_count:desc"}, {"changefreq", "daily"}, {"priority", 0.6}})
	sm.Add(stm.URL{{"loc", "/search?sort=recent_ask:desc"}, {"changefreq", "daily"}, {"priority", 0.6}})
	sm.Add(stm.URL{{"loc", "/about"}})
	sm.Add(stm.URL{{"loc", "/faq"}})
	sm.Add(stm.URL{{"loc", "/privacy"}})

	// Catalog listings locations.
	catalogs, _, _ := marketSvc.Catalog(core.FindOpts{})
	for _, cc := range catalogs {
		sm.Add(stm.URL{{"loc", "/item/" + cc.Slug}, {"changefreq", "daily"}, {"priority", 0.7}})
	}

	// User profile locations.
	users, _ := userSvc.Users(core.FindOpts{})
	for _, uu := range users {
		sm.Add(stm.URL{{"loc", "/user/" + uu.SteamID}, {"changefreq", "daily"}, {"priority", 0.7}})
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(sm.XMLContent())
		return
	}
}
