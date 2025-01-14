package main

import (
	"bbb-voting/prodution-frontend/controller"
	_ "bbb-voting/prodution-frontend/docs"
	postgresqldatamapper "bbb-voting/voting-commons/data-layer/postgresql"
	redisdatamapper "bbb-voting/voting-commons/data-layer/redis"
	"bbb-voting/voting-commons/domain"
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
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
	redisClient := getRedisClient(os.Getenv("REDIS_URL"))

	var participantRepository domain.ParticipantRepository = postgresqldatamapper.NewParticipantDataMapper(
		postgresqlConnector,
	)
	participantRepository = redisdatamapper.DecorateParticipantDataMapperWithRedis(participantRepository, *redisClient)

	frontendController := controller.NewFrontendController(
		participantRepository,
		context, templates,
	)
	http.HandleFunc("/", frontendController.GetPage)
	http.HandleFunc("/votes/totals/thorough", frontendController.GetThoroughTotals)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFiles))))
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	log.Println("Server is running on http://localhost:8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}

func getRedisClient(url string) *redis.Client {
	opts, err := redis.ParseURL(url)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}
	return redis.NewClient(opts)
}
