package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kudarap/dotagiftx/dota2"
)

func handleTreasureList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dota2ContentCacheControl(w)
		respondOK(w, dota2.AllTreasures)
	}
}

func handleHeroList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dota2ContentCacheControl(w)
		respondOK(w, dota2.AllHeroes)
	}
}

func dota2ContentCacheControl(w http.ResponseWriter) {
	const age = time.Hour * 24 * 30
	cc := fmt.Sprintf("max-age=%d, public", int(age.Seconds()))
	w.Header().Add("Cache-Control", cc)
}
