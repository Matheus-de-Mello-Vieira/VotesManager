package main

import (
	"bbb-voting/voters-frontend/controller"
	_ "bbb-voting/voters-frontend/docs"
	kafkadatamapper "bbb-voting/voting-commons/data-layer/kafka"
	postgresqldatamapper "bbb-voting/voting-commons/data-layer/postgresql"
	redisdatamapper "bbb-voting/voting-commons/data-layer/redis"
	"bbb-voting/voting-commons/domain"
	"context"
	"embed"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
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

	ctx := context.Background()
	postgresqlConnector := postgresqldatamapper.NewPostgresqlConnector(os.Getenv("POSTGRESQL_URI"))
	redisClient := getRedisClient(os.Getenv("REDIS_URL"))

	var participantRepository domain.ParticipantRepository = postgresqldatamapper.NewParticipantDataMapper(
		postgresqlConnector,
	)
	var err error;
	participantRepository, err = redisdatamapper.DecorateParticipantDataMapperWithRedis(participantRepository, *redisClient, ctx)
	if err != nil {
		log.Fatalf("Faled to load cache: %s", err)
	}

	var voteRepository domain.VoteRepository = kafkadatamapper.NewVoteDataMapper(
		[]string{os.Getenv("KAFKA_URI")}, "votes", 30,
	)
	voteRepository = redisdatamapper.DecorateVoteDataMapperWithRedis(voteRepository, *redisClient)

	frontendController := controller.NewFrontendController(
		participantRepository,
		voteRepository,
		ctx, templates,
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

func getRedisClient(url string) *redis.Client {
	opts, err := redis.ParseURL(url)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}
	return redis.NewClient(opts)
}
