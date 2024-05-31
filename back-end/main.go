package main

import (
	_ "embed"
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var (
	router *mux.Router = mux.NewRouter()
	port   string
)

func handleUserRoutes() {
	router.HandleFunc("/user/", GetUserHandler).Methods("GET")
	router.HandleFunc("/user/post", PutUserHandler).Methods("POST")
	router.HandleFunc("/user/del", DelUserHandler).Methods("DELETE")
	router.HandleFunc("/user/list", GetUsersListHandler).Methods("GET")
}

func handleDateRoutes() {
	router.HandleFunc("/date/list", GetDateListHandler).Methods("GET")
	router.HandleFunc("/date/post", PutDateHandler).Methods("POST")
	router.HandleFunc("/date/del", DelDateHandler).Methods("DELETE")
}

func main() {
	flag.StringVar(&port, "port", ":8080", "The port the server will use to listen to requests")
	flag.Parse()

	handleUserRoutes()
	handleDateRoutes()

	handler := cors.Default().Handler(router)
	log.Print("server will be listening on port ", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
