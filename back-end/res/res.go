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
		log.Println("res.Err - json.Marshal ", err)
		return
	}
	w.WriteHeader(status)
	if _, err = w.Write(json); err != nil {
		log.Println("res.Err - w.Write ", err)
	}

}

func Ok(w http.ResponseWriter, item interface{}) {
	res := Response{
		Base: BaseAPI{Ok: true},
		Res:  item,
	}

	json, err := json.Marshal(res)
	if err != nil {
		log.Println("res.Ok - json.Marshal ", err)
		return
	}
	if _, err = w.Write(json); err != nil {
		log.Println("res.Ok - w.Write ", err)
	}
}
