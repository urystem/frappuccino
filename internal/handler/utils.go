package handler

import (
	"encoding/json"
	"net/http"
	"regexp"
)

func writeHttp(w http.ResponseWriter, code int, where, errOrMes string) error {
	key := "error"
	if code < 300 {
		key = "message"
	}
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(map[string]string{key: where + " : " + errOrMes})
}

func bodyJsonStruct(w http.ResponseWriter, someThing any) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(someThing)
}

func checkName(name string) bool {
	return !regexp.MustCompile(`^[ \w+]{1,128}$`).MatchString(name) || regexp.MustCompile("  ").MatchString(name) || name[0] == ' ' || name[len(name)-1] == ' '
}
