package http

import (
	"net/http"

	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"
	"github.com/kudarap/dotagiftx/core"
)

func buildSitemap(items []core.Item, users []core.User) *stm.Sitemap {
	sitemap := stm.NewSitemap(1)
	//sitemap.SetVerbose(false)
	sitemap.SetDefaultHost("https://dotagiftx.com")
	sitemap.Create()

	// Add static pages locations.
	sitemap.Add(stm.URL{{"loc", "/"}, {"changefreq", "daily"}, {"priority", 0.8}})
	sitemap.Add(stm.URL{{"loc", "/search"}, {"changefreq", "daily"}, {"priority", 0.6}})
	sitemap.Add(stm.URL{{"loc", "/search?sort=" + queryFlagRecentItems}, {"changefreq", "daily"}, {"priority", 0.6}})
	sitemap.Add(stm.URL{{"loc", "/search?sort=" + queryFlagPopularItems}, {"changefreq", "daily"}, {"priority", 0.6}})
	sitemap.Add(stm.URL{{"loc", "/about"}})
	sitemap.Add(stm.URL{{"loc", "/faq"}})
	sitemap.Add(stm.URL{{"loc", "/privacy"}})
	sitemap.Add(stm.URL{{"loc", "/login"}})

	// Add item slug locations.
	origins := map[string]struct{}{}
	heroes := map[string]struct{}{}
	for _, ii := range items {
		sitemap.Add(stm.URL{{"loc", "/" + ii.Slug}, {"changefreq", "daily"}, {"priority", 0.7}})
		origins[ii.Origin] = struct{}{}
		heroes[ii.Hero] = struct{}{}
	}
	// Add item origin and heroes.
	for i, _ := range origins {
		sitemap.Add(stm.URL{{"loc", "/search?q=" + i}})
	}
	for i, _ := range heroes {
		sitemap.Add(stm.URL{{"loc", "/search?q=" + i}})
	}

	// Add user profile locations.
	for _, uu := range users {
		sitemap.Add(stm.URL{{"loc", "/profiles/" + uu.SteamID}, {"changefreq", "daily"}, {"priority", 0.6}})
	}

	return sitemap
}

func handleSitemap(itemSvc core.ItemService, userSvc core.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, _, _ := itemSvc.Items(core.FindOpts{})
		users, _ := userSvc.Users(core.FindOpts{})

		w.Header().Set("content-type", "text/xml")
		w.Write(buildSitemap(items, users).XMLContent())
		return
	}
}
