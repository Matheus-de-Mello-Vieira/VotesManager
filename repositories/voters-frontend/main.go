package main

import (
	"bbb-voting/voters-frontend/controller"
	"bbb-voting/voting-commons/data-layer/postgresql"
	"context"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	
	context := context.Background()
	postgresqlConnector := postgresqldatamapper.NewPostgresqlConnector(os.Getenv("postgresql_uri"))
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

	log.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
