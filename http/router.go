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
			r.Get("/{id}", handleImage(s.imageSvc))
			r.Get("/{id}/{w}x{h}", handleImageThumbnail(s.imageSvc))
		})
		r.Get("/users/{id}", handlePublicProfile(s.userSvc))
		r.Route("/items", func(r chi.Router) {
			r.Get("/", handleItemList(s.itemSvc))
			r.Get("/{id}", handleItemDetail(s.itemSvc))
		})
	})
}

func (s *Server) privateRouter(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(s.authorizer)
		r.Route("/my", func(r chi.Router) {
			r.Get("/profile", handleProfile(s.userSvc))
			r.Route("/sells", func(r chi.Router) {
				r.Get("/", handleSellList(s.sellSvc))
				r.Post("/", handleSellCreate(s.sellSvc))
				r.Get("/{id}", handleSellDetail(s.sellSvc))
				r.Patch("/{id}", handleSellUpdate(s.sellSvc))
			})
			r.Post("/images", handleImageUpload(s.imageSvc))
		})
		r.Post("/items", handleItemCreate(s.itemSvc))
	})
}
