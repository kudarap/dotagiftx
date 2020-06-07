package main

import (
	"log"
	"net/http"

	"github.com/kudarap/dota2giftables/steam"
)

var client *steam.Client

func main() {
	c, err := steam.New(steam.Config{"B4F4D3D11EDFD1E208378B272971A5AB"})
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

	sp, err := client.Verify(r)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(sp)
	w.Write([]byte(sp.Name))
}
