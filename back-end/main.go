package main

import (
	_ "embed"
	"log"
	"net/http"
	"remindal/res"

	"github.com/gorilla/mux"
)

var (
	//go:embed token
	embed_test string
	port       string = ":8080"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	if id := r.URL.Query().Get("id"); id == "" {
		res.Err(w, http.StatusBadRequest)
		return
	}
	res.Ok(w)
}

func main() {
	log.Print(embed_test)
	router := mux.NewRouter()

	router.HandleFunc("/", handleHome)
	router.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		res.Err(w, http.StatusBadRequest)
	})

	log.Print("server will be listening on port", port)
	log.Fatal(http.ListenAndServe(port, router))
}
