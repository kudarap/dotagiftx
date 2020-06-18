package http

import (
	"net/http"

	"github.com/kudarap/dota2giftables/core"
)

func handleSellList(svc core.SellService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respondOK(w, newMsg("not implemented"))
	}
}

func handleSellDetail(svc core.SellService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respondOK(w, newMsg("not implemented"))
	}
}

func handleSellCreate(svc core.SellService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respondOK(w, newMsg("not implemented"))
	}
}

func handleSellUpdate(svc core.SellService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respondOK(w, newMsg("not implemented"))
	}
}
