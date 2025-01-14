package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func loadBody(responseWriter http.ResponseWriter, request *http.Request, contentBody any) {
	bytesBody, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(responseWriter, "Error when read", http.StatusBadRequest)
	}

	err = json.Unmarshal(bytesBody, &contentBody)
	if err != nil {
		http.Error(responseWriter, "Error when read", http.StatusMethodNotAllowed)
	}
}

func handleInternalServerError(responseWriter http.ResponseWriter, err error) {
	http.Error(responseWriter, "Internal Server Error", http.StatusInternalServerError)
	log.Fatal(err)
}
