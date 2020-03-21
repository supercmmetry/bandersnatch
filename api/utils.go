package api

import (
	"encoding/json"
	"net/http"
)

func RespWrap(w http.ResponseWriter, statusCode int, msg string) {
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"msg": msg})
}

func JsonifyHeader(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
}

func Wrap(w http.ResponseWriter, v map[string]interface{}) {
	JsonifyHeader(w)
	_ = json.NewEncoder(w).Encode(v)
}


