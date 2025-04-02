package handler

import (
	"encoding/json"
	"net/http"
)

func writeHttp(w http.ResponseWriter, code int, where, errOrMes string) error {
	key := "error"
	if code < 300 {
		key = "message"
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(map[string]string{key: where + " : " + errOrMes})
}

func bodyJsonStruct(w http.ResponseWriter, someThing any, code int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(someThing)
}
