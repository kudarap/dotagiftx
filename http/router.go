package http

import "github.com/go-chi/chi"

func (s *Server) publicRouter(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Get("/", handleInfo(s.version))

		r.Route("/auth", func(r chi.Router) {
			r.Get("/steam", handleAuthTwitter(s.authSvc))
			r.Post("/renew", handleAuthRenew(s.authSvc))
			r.Post("/revoke", handleAuthRevoke(s.authSvc))
		})

		r.Route("/images", func(r chi.Router) {
			r.Get("/{id}", handleImage(s.imageSvc))
			r.Get("/{id}/{w}x{h}", handleImageThumbnail(s.imageSvc))
		})

		r.Route("/items", func(r chi.Router) {
			r.Get("/", handleItemList(s.itemSvc))
			r.Get("/{id}", handleItemDetail(s.itemSvc))
		})

		r.Get("/users/{id}", handlePublicProfile(s.userSvc))
	})
}

func (s *Server) privateRouter(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(s.authorizer)
		r.Route("/me", func(r chi.Router) {
			r.Get("/", handleProfile(s.userSvc))
			r.Post("/images", handleImageUpload(s.imageSvc))
			r.Route("/items", func(r chi.Router) {
				r.Get("/", handleItemList(s.itemSvc))
				r.Post("/", handleItemCreate(s.itemSvc))
				r.Get("/{id}", handleItemDetail(s.itemSvc))
				r.Patch("/{id}", handleItemUpdate(s.itemSvc))
			})
		})
	})
}
