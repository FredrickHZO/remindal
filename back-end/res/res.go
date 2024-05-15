package res

import (
	"encoding/json"
	"log"
	"net/http"
	remerr "remindal/errors"
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
func Err(w http.ResponseWriter, err error) {
	res := ResponseAPI{Ok: false, Message: err.Error()}

	json, err := json.Marshal(res)
	if err != nil {
		log.Println("res.Err - json.Marshal ", err)
		return
	}

	status := statusError(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err = w.Write(json); err != nil {
		log.Println("res.Err - w.Write ", err)
	}
}

/*
Returns the appropriate HTTP status code based on the given error.
*/
func statusError(err error) int {
	var status int
	if err == remerr.ErrInternalServerError {
		status = http.StatusInternalServerError
	} else {
		status = http.StatusBadRequest
	}
	return status
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
