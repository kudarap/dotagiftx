package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"
	dgx "github.com/kudarap/dotagiftx"
)

func buildSitemap(items []dgx.Item, users []dgx.User, vanities []string) *stm.Sitemap {
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
	sitemap.Add(stm.URL{{"loc", "/guides"}})
	sitemap.Add(stm.URL{{"loc", "/rules"}})

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

func handleSitemap(itemSvc dgx.ItemService, userSvc dgx.UserService, cache dgx.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, _, _ := itemSvc.Items(dgx.FindOpts{})
		users, _ := userSvc.Users(dgx.FindOpts{
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
