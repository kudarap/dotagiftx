package http

import (
	"net/http"
	"strings"

	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"
	"github.com/kudarap/dotagiftx/core"
)

func buildSitemap(items []core.Item, users []core.User, vanities []string) *stm.Sitemap {
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
	sitemap.Add(stm.URL{{"loc", "/faqs"}})
	sitemap.Add(stm.URL{{"loc", "/privacy"}})
	sitemap.Add(stm.URL{{"loc", "/login"}})
	sitemap.Add(stm.URL{{"loc", "/donate"}})
	sitemap.Add(stm.URL{{"loc", "/middlemen"}})

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
		sitemap.Add(stm.URL{{"loc", "/search?origin=" + i}})
	}
	for i, _ := range heroes {
		sitemap.Add(stm.URL{{"loc", "/search?hero=" + i}})
	}

	// Add user profile locations.
	for _, uu := range users {
		sitemap.Add(stm.URL{{"loc", "/profiles/" + uu.SteamID}, {"changefreq", "daily"}, {"priority", 0.6}})
	}

	// Add user vanity urls locations.
	for _, v := range vanities {
		sitemap.Add(stm.URL{{"loc", "/id/" + v}, {"changefreq", "daily"}, {"priority", 0.6}})
	}

	return sitemap
}

const vanityPrefix = "https://steamcommunity.com/id/"

func handleSitemap(itemSvc core.ItemService, userSvc core.UserService, steam core.SteamClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, _, _ := itemSvc.Items(core.FindOpts{})
		users, _ := userSvc.Users(core.FindOpts{})

		var vanities []string
		for _, u := range users {
			//sp, _ := steam.Player(u.SteamID)
			//if sp == nil || sp.URL == "" {
			//	continue
			//}
			sp := u

			// Not a custom url.
			if !strings.HasPrefix(sp.URL, vanityPrefix) {
				continue
			}

			vanities = append(vanities, strings.TrimPrefix(sp.URL, vanityPrefix))
		}

		w.Header().Set("content-type", "text/xml")
		w.Write(buildSitemap(items, users, vanities).XMLContent())
		return
	}
}
