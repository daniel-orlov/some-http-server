package transport

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func renderErrorResponse(w http.ResponseWriter, msg string, status int) {
	renderResponse(w, errorResponse{Error: msg}, status)
}

func renderResponse(w http.ResponseWriter, res interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")

	content, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)

	if _, err = w.Write(content); err != nil {
		// TODO handle error
	}
}
