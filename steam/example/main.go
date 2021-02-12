package main

import (
	"log"
	"net/http"

	"github.com/kudarap/dotagiftx/steam"
)

var client *steam.Client

func main() {
	c, err := steam.New(steam.Config{Key: "STEAM_WEB_API_KEY"}, nil)
	if err != nil {
		log.Fatalln(err)
	}
	client = c

	http.HandleFunc("/login", loginHandler)
	http.ListenAndServe(":9000", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("openid.mode") == "" {
		url, err := client.AuthorizeURL(r)
		if err != nil {
			log.Fatalln(err)
		}

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}

	sp, err := client.Authenticate(r)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(sp)
	w.Write([]byte(sp.Name))
}
