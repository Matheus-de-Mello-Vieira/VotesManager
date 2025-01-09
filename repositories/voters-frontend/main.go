package main

import (
	"log"
	"net/http"
	"bbb-voting/voters-frontend/controller"
)

func main() {
	http.HandleFunc("/", controller.IndexHandler)

	log.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
