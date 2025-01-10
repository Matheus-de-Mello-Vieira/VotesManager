package main

import (
	"bbb-voting/prodution-frontend/controller"
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
		context, templates,
	)
	http.HandleFunc("/", frontendController.GetPage)
	http.HandleFunc("/votes/thorough", frontendController.GetThoroughTotals)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFiles))))

	log.Println("Server is running on http://localhost:8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
