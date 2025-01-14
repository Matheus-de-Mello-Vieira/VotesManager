package main

import (
	"bbb-voting/voters-frontend/controller"
	_ "bbb-voting/voters-frontend/docs"
	kafkadatamapper "bbb-voting/voting-commons/data-layer/kafka"
	postgresqldatamapper "bbb-voting/voting-commons/data-layer/postgresql"
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	httpSwagger "github.com/swaggo/http-swagger"
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
	http.HandleFunc("/pages/totals/rough", frontendController.LoadRoughTotalPage)

	http.HandleFunc("/votes", frontendController.PostVoteHandler)
	http.HandleFunc("/participants", frontendController.GetParticipantsHandler)
	http.HandleFunc("/votes/totals/rough", frontendController.GetVotesRoughTotalsHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFiles))))
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	log.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
