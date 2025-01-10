package main

import (
	"bbb-voting/voters-frontend/controller"
	kafkadatamapper "bbb-voting/voting-commons/data-layer/kafka"
	postgresqldatamapper "bbb-voting/voting-commons/data-layer/postgresql"
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
)

//go:embed view/static/*
var staticFilesFull embed.FS

//go:embed view/templates/*
var templatesFull embed.FS

func main() {
	var templates, _ = fs.Sub(templatesFull, "view/templates")
	var staticFiles, _ = fs.Sub(staticFilesFull, "view/static")

	context := context.Background()
	postgresqlConnector := postgresqldatamapper.NewPostgresqlConnector(os.Getenv("POSTGRESQL_URI"))
	frontendController := controller.NewFrontendController(
		postgresqldatamapper.NewParticipantDataMapper(
			postgresqlConnector,
		),
		kafkadatamapper.NewVoteDataMapper(
			[]string{os.Getenv("KAFKA_URI")}, "votes", 30,
		),
		context, templates,
	)
	http.HandleFunc("/", frontendController.IndexHandler)
	http.HandleFunc("/votes", frontendController.VoteCastingHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFiles))))

	log.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
