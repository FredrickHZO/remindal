package res

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func Err(w http.ResponseWriter, status int) {
	res := Response{Message: "error"}

	json, err := json.Marshal(res)
	if err != nil {
		log.Fatal("eres - json.Marshal")
	}
	w.WriteHeader(status)
	w.Write(json)
}

func Ok(w http.ResponseWriter) {
	res := Response{Message: "ok"}

	json, err := json.Marshal(res)
	if err != nil {
		log.Fatal("okres - json.Marshal")
	}
	w.Write(json)
}
