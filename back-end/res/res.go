package res

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
and writes the header with the specified HTTP error status.
*/
func Err(w http.ResponseWriter, err error, status int) {
	res := ResponseAPI{Ok: false, Message: err.Error()}

	json, err := json.Marshal(res)
	if err != nil {
		log.Println("res.Err - json.Marshal ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err = w.Write(json); err != nil {
		log.Println("res.Err - w.Write ", err)
	}
}

/*
Sends a successful response to the client and writes data
- if any - as part of the HTTP response body.
*/
func Ok(w http.ResponseWriter, item any) {
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
