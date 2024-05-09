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

type Response struct {
	Base BaseAPI
	Res  interface{} `json:"res,omitempty"`
}

func Err(w http.ResponseWriter, err error, status int) {
	res := BaseAPI{Ok: false, Message: err.Error()}

	json, err := json.Marshal(res)
	if err != nil {
		log.Fatal("Err - json.Marshal")
		return
	}
	w.WriteHeader(status)
	w.Write(json)
}

func Ok(w http.ResponseWriter, item interface{}) {
	res := Response{
		Base: BaseAPI{Ok: true},
		Res:  item,
	}

	json, err := json.Marshal(res)
	if err != nil {
		log.Fatal("okres - json.Marshal")
		return
	}
	w.Write(json)
}
