package main

import (
	_ "embed"
	"flag"
	"log"
	"net/http"
	"remindal/routes"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var (
	router *mux.Router = mux.NewRouter()
	port   string
)

func handleTestRoutes() {
	router.HandleFunc("/user", routes.GetUserHandler).Methods("GET")
	router.HandleFunc("/user", routes.PutUserHandler).Methods("POST")
	router.HandleFunc("/user", routes.DelUserHandler).Methods("DELETE")
	router.HandleFunc("/user/list", routes.GetUsersListHandler).Methods("GET")
}

func main() {
	flag.StringVar(&port, "port", ":8080", "The port the server will use to listen to requests")
	flag.Parse()

	handleTestRoutes()

	handler := cors.Default().Handler(router)
	log.Print("server will be listening on port ", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
