package main

import (
	_ "embed"
	"log"
	"net/http"
	"remindal/routes"

	"github.com/gorilla/mux"
)

var (
	router *mux.Router = mux.NewRouter()
	port   string      = ":8080"
)

func handleTestRoutes() {
	router.HandleFunc("/user", routes.GetUserHandle).Methods("GET")
	router.HandleFunc("/user", routes.PutUserHandle).Methods("PUT")
	router.HandleFunc("/user/list", routes.GetUsersListHandle).Methods("GET")
}

func main() {
	handleTestRoutes()

	log.Print("server will be listening on port", port)
	log.Fatal(http.ListenAndServe(port, router))
}
