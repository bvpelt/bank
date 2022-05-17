package test

import (
	"net/http"

	sw "github.com/bvpelt/bank/swagger"

	log "github.com/sirupsen/logrus"
)

func Server() {
	/*
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		})
	*/

	log.Println("Listening on localhost:8080")

	router := sw.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
