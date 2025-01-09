package main

import (
	"log"
	"net/http"
	"bbb-voting/voters-frontend/controller"
)

func main() {
	http.HandleFunc("/", controller.IndexHandler)
	http.HandleFunc("/vote", controller.VoteCastingHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("voters-frontend/view/static"))))

	log.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
