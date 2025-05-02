package api

import (
	"encoding/json"
	"net/http"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

type Err struct {
	Err string `json:"err"`
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, Err{Err: err.Error()})
		}
	}
}
