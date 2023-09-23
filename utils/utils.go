package utils

import (
	"encoding/json"
	"net/http"
)



func ResponseWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	out, _ := json.Marshal(payload)
	w.Write(out)
}


func ResponseWithError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	eMsg := struct {
		Error string `json:"error"`
	}{ Error: msg }
	out, _ := json.Marshal(eMsg)
	w.Write(out)
}
