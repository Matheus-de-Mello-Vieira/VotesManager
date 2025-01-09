package main

import (
	"bbb-voting/prodution-frontend/controller"
	postgresqldatamapper "bbb-voting/voting-commons/data-layer/postgresql"
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
		context,
	)
	http.HandleFunc("/votes/detailed", frontendController.GetThoroughTotals)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("voters-frontend/view/static"))))

	log.Println("Server is running on http://localhost:8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
