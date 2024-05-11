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
	router.HandleFunc("/user", routes.GetUserHandler).Methods("GET")
	router.HandleFunc("/user", routes.PutUserHandler).Methods("PUT")
	router.HandleFunc("/user/list", routes.GetUsersListHandler).Methods("GET")
}

func main() {
	handleTestRoutes()

	log.Print("server will be listening on port", port)
	log.Fatal(http.ListenAndServe(port, router))
}
