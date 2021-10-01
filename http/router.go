package http

import "github.com/go-chi/chi"

func (s *Server) publicRouter(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Get("/", handleInfo(s.version))
		r.Route("/auth", func(r chi.Router) {
			r.Get("/steam", handleAuthSteam(s.authSvc))
			r.Post("/renew", handleAuthRenew(s.authSvc))
			r.Post("/revoke", handleAuthRevoke(s.authSvc))
		})
		r.Route("/images", func(r chi.Router) {
			r.Get("/{w}x{h}/{id}", handleImageThumbnail(s.imageSvc))
			r.Get("/{id}", handleImage(s.imageSvc))
		})
		r.Route("/items", func(r chi.Router) {
			r.Get("/", handleItemList(s.itemSvc, s.trackSvc, s.cache, s.logger))
			r.Get("/{id}", handleItemDetail(s.itemSvc, s.cache, s.logger))
		})
		r.Route("/markets", func(r chi.Router) {
			r.Get("/", handleMarketList(s.marketSvc, s.trackSvc, s.cache, s.logger))
			r.Get("/{id}", handleMarketDetail(s.marketSvc, s.cache, s.logger))
		})
		r.Get("/catalogs_trend", handleMarketCatalogTrendList(s.marketSvc, s.cache, s.logger))
		r.Get("/catalogs", handleMarketCatalogList(s.marketSvc, s.trackSvc, s.cache, s.logger))
		r.Get("/catalogs/{slug}", handleMarketCatalogDetail(s.marketSvc, s.cache, s.logger))
		r.Get("/users/{id}", handlePublicProfile(s.userSvc, s.cache))
		r.Get("/t", handleTracker(s.trackSvc, s.logger))
		r.Get("/sitemap.xml", handleSitemap(s.itemSvc, s.userSvc, s.steam))
		r.Get("/stats/market_summary", handleStatsMarketSummary(s.statsSvc, s.cache))
		r.Get("/stats/top_origins", handleStatsTopOrigins(s.itemSvc, s.cache))
		r.Get("/stats/top_heroes", handleStatsTopHeroes(s.itemSvc, s.cache))
		r.Get("/graph/market_sales", handleGraphMarketSales(s.statsSvc, s.cache))
		r.Route("/reports", func(r chi.Router) {
			r.Get("/", handleReportList(s.reportSvc))
			r.Get("/{id}", handleReportDetail(s.reportSvc))
		})
		r.Get("/vanity/{id}", handleVanityProfile(s.userSvc, s.steam, s.cache))
		r.Get("/blacklists", handleBlacklisted(s.userSvc, s.cache))
	})
}

func (s *Server) privateRouter(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(s.authorizer)
		r.Route("/my", func(r chi.Router) {
			r.Get("/profile", handleProfile(s.userSvc, s.cache))
			r.Route("/markets", func(r chi.Router) {
				r.Get("/", handleMarketList(s.marketSvc, s.trackSvc, s.cache, s.logger))
				r.Post("/", handleMarketCreate(s.marketSvc, s.cache))
				r.Get("/{id}", handleMarketDetail(s.marketSvc, s.cache, s.logger))
				r.Patch("/{id}", handleMarketUpdate(s.marketSvc, s.cache))
			})
		})
		r.Post("/items", handleItemCreate(s.itemSvc, s.cache))
		r.Post("/items_import", handleItemImport(s.itemSvc, s.cache))
		r.Post("/images", handleImageUpload(s.imageSvc))
		r.Post("/reports", handleReportCreate(s.reportSvc))
		r.Post("/hammer/ban", handleHammerBan(s.hammerSvc))
		r.Post("/hammer/suspend", handleHammerSuspend(s.hammerSvc))
		r.Post("/hammer/lift", handleHammerLift(s.hammerSvc))
	})
}
