package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"
	"github.com/kudarap/dotagiftx"
)

func buildSitemap(items []dotagiftx.Item, users []dotagiftx.User, vanities []string) *stm.Sitemap {
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
	sitemap.Add(stm.URL{{"loc", "/guides"}})
	sitemap.Add(stm.URL{{"loc", "/rules"}})
	sitemap.Add(stm.URL{{"loc", "/plus"}})
	sitemap.Add(stm.URL{{"loc", "/updates"}})
	sitemap.Add(stm.URL{{"loc", "/treasures"}})
	sitemap.Add(stm.URL{{"loc", "/moderators"}})
	sitemap.Add(stm.URL{{"loc", "/middleman"}})

	// Add item slug locations.
	origins := map[string]struct{}{}
	heroes := map[string]struct{}{}
	for _, ii := range items {
		sitemap.Add(stm.URL{{"loc", "/" + ii.Slug}, {"changefreq", "daily"}, {"priority", 0.7}})
		origins[ii.Origin] = struct{}{}
		heroes[ii.Hero] = struct{}{}
	}
	// Add item origin and heroes.
	for i := range origins {
		sitemap.Add(stm.URL{{"loc", "/search?origin=" + i}})
	}
	for i := range heroes {
		sitemap.Add(stm.URL{{"loc", "/search?hero=" + i}})
	}

	// Add user profile locations.
	//for _, uu := range users {
	//	sitemap.Add(stm.URL{{"loc", "/profiles/" + uu.SteamID}, {"changefreq", "monthly"}, {"priority", 0.6}})
	//}
	// Add user vanity urls locations.
	//for _, v := range vanities {
	//	sitemap.Add(stm.URL{{"loc", "/id/" + v}, {"changefreq", "monthly"}, {"priority", 0.6}})
	//}

	return sitemap
}

const (
	vanityPrefix = "https://steamcommunity.com/id/"

	sitemapCacheKey  = "sitemap"
	sitemapCacheExpr = time.Hour
)

func handleSitemap(itemSvc dotagiftx.ItemService, userSvc dotagiftx.UserService, cache dotagiftx.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, _, _ := itemSvc.Items(dotagiftx.FindOpts{})
		users, _ := userSvc.Users(dotagiftx.FindOpts{
			Limit: 0,
		})

		var vanities []string
		for _, u := range users {
			sp := u
			if !strings.HasPrefix(sp.URL, vanityPrefix) {
				continue
			}

			vanities = append(vanities, strings.TrimPrefix(sp.URL, vanityPrefix))
		}

		sm := buildSitemap(items, users, vanities).XMLContent()
		w.Header().Set("content-type", "text/xml")
		if _, err := w.Write(sm); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
}
