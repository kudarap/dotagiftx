package http

import (
	"net/http"

	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"
	"github.com/kudarap/dotagiftx/core"
)

func buildSitemap(items []core.Item, users []core.User) *stm.Sitemap {
	sm := stm.NewSitemap(1)
	//sm.SetVerbose(false)
	sm.SetDefaultHost("https://dotagiftx.com")
	sm.Create()

	// Static pages locations.
	sm.Add(stm.URL{{"loc", "/"}, {"changefreq", "daily"}, {"priority", 0.8}})
	sm.Add(stm.URL{{"loc", "/search"}, {"changefreq", "daily"}, {"priority", 0.6}})
	sm.Add(stm.URL{{"loc", "/search?sort=" + queryFlagRecentItems}, {"changefreq", "daily"}, {"priority", 0.6}})
	sm.Add(stm.URL{{"loc", "/search?sort=" + queryFlagPopularItems}, {"changefreq", "daily"}, {"priority", 0.6}})
	sm.Add(stm.URL{{"loc", "/about"}})
	sm.Add(stm.URL{{"loc", "/faq"}})
	sm.Add(stm.URL{{"loc", "/privacy"}})
	sm.Add(stm.URL{{"loc", "/login"}})

	// Catalog listings locations.
	for _, ii := range items {
		sm.Add(stm.URL{{"loc", "/item/" + ii.Slug}, {"changefreq", "daily"}, {"priority", 0.7}})
	}

	// User profile locations.
	for _, uu := range users {
		sm.Add(stm.URL{{"loc", "/user/" + uu.SteamID}, {"changefreq", "daily"}, {"priority", 0.7}})
	}

	return sm
}

func handleSitemap(itemSvc core.ItemService, userSvc core.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, _, _ := itemSvc.Items(core.FindOpts{})
		users, _ := userSvc.Users(core.FindOpts{})
		w.Write(buildSitemap(items, users).XMLContent())
		return
	}
}
