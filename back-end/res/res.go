package res

import (
	"encoding/json"
	"log"
	"net/http"
)

type BaseAPI struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
}

type ResponseAPI struct {
	Base BaseAPI
	Res  interface{} `json:"res,omitempty"`
}

/*
Sends to the client an error response with a description
and writes the header with the specified HTTP error status
*/
func Err(w http.ResponseWriter, err error, status int) {
	res := BaseAPI{Ok: false, Message: err.Error()}

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
Sends the client a successful response and writes data - if any - as
part of the HTTP response in the body.

Automatically sets the HTTP status to 200.
*/
func Ok(w http.ResponseWriter, item any) {
	res := ResponseAPI{
		Base: BaseAPI{Ok: true},
		Res:  item,
	}

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
