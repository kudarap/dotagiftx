package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kudarap/dotagiftx/dota2"
)

func handleTreasureList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dota2ContentCacheControl(w)
		respondOK(w, dota2.AllTreasures)
	}
}

func handleTreasureDetail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := dota2.TreasureDetail(chi.URLParam(r, "slug"))
		if err != nil {
			respondError(w, err)
			return
		}
		respondOK(w, data)
	}
}

func handleHeroList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dota2ContentCacheControl(w)
		respondOK(w, dota2.AllHeroes)
	}
}

func handleHeroDetail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := dota2.HeroDetail(chi.URLParam(r, "id"))
		if err != nil {
			respondError(w, err)
			return
		}
		respondOK(w, data)
	}
}

func dota2ContentCacheControl(w http.ResponseWriter) {
	const age = time.Hour * 24 * 30
	cc := fmt.Sprintf("max-age=%d, public", int(age.Seconds()))
	w.Header().Add("Cache-Control", cc)
}
