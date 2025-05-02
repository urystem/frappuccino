package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func writeHttp(w http.ResponseWriter, code int, where, errOrMes string) {
	key := "error"
	if code < 300 {
		key = "message"
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(map[string]string{key: where + " : " + errOrMes})
	if err != nil {
		slog.Error("cannot write to w this message", key, where+":"+errOrMes)
	}
}

func bodyJsonStruct(w http.ResponseWriter, someThing any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(someThing)
	if err != nil {
		slog.Error("bodyJsonStruct error:", "", err)
	}
}
