package main

import (
	"bbb-voting/voters-frontend/controller"
	postgresqldatamapper "bbb-voting/voting-commons/data-layer/postgresql"
	"context"
	"log"
	"net/http"
	"os"
)

func main() {
	context := context.Background()
	postgresqlConnector := postgresqldatamapper.NewPostgresqlConnector(os.Getenv("POSTGRESQL_URI"))
	frontendController := controller.NewFrontendController(
		postgresqldatamapper.NewParticipantDataMapper(
			postgresqlConnector,
		),
		postgresqldatamapper.NewVoteDataMapper(
			postgresqlConnector,
		),
		context,
	)
	http.HandleFunc("/", frontendController.IndexHandler)
	http.HandleFunc("/votes", frontendController.VoteCastingHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("voters-frontend/view/static"))))
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("voters-frontend/view/templates"))))

	log.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
