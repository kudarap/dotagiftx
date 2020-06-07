package http

import (
	"net/http"

	"github.com/kudarap/dota2giftables/core"
)

func handleItemList(svc core.ItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respondOK(w, newMsg("not implemented"))
	}
}

func handleItemDetail(svc core.ItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respondOK(w, newMsg("not implemented"))
	}
}

func handleItemCreate(svc core.ItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respondOK(w, newMsg("not implemented"))
	}
}

func handleItemUpdate(svc core.ItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respondOK(w, newMsg("not implemented"))
	}
}
