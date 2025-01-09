package main

import (
	"bbb-voting/voters-frontend/controller"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", controller.IndexHandler)
	http.HandleFunc("/votes", controller.VoteCastingHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("voters-frontend/view/static"))))

	log.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
