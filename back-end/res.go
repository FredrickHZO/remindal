package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResponseAPI struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Res     any    `json:"res,omitempty"`
}

/*
Sends an error response to the client with a description
and automatically detects the appropriate HTTP error status,
writing it to the header.
*/
func Eres(w http.ResponseWriter, se *HttpError) {
	res := ResponseAPI{Ok: false, Message: se.Error()}

	json, err := json.Marshal(res)
	if err != nil {
		log.Println("res.Err - json.Marshal ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(se.status)
	if _, err = w.Write(json); err != nil {
		log.Println("res.Err - w.Write ", err)
	}
}

/*
Sends a successful response to the client and writes data
- if any - as part of the HTTP response body.
*/
func Okres(w http.ResponseWriter, item any) {
	res := ResponseAPI{Ok: true, Res: item}

	json, err := json.Marshal(res)
	if err != nil {
		log.Println("res.Ok - json.Marshal ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(json); err != nil {
		log.Println("res.Ok - w.Write ", err)
	}
}

func Test() {
	var jsn []byte

	var res ResponseAPI
	err := json.Unmarshal(jsn, &res)
	if err != nil {
		log.Println("ha sciut")
	}
}
