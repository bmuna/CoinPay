package server

import (
	"encoding/json"
	"log"
	"net/http"
)


func responseWithJSON(w http.ResponseWriter, code int, payload interface{}){
	data, err := json.Marshal(payload)

	if err != nil {
		log.Printf("failed to marshal a json %v", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)


}