package main

import (
	"bbb-voting/voters-frontend/controller"
	kafkadatamapper "bbb-voting/voting-commons/data-layer/kafka"
	postgresqldatamapper "bbb-voting/voting-commons/data-layer/postgresql"
	redisdatamapper "bbb-voting/voting-commons/data-layer/redis"
	"bbb-voting/voting-commons/domain"
	usercases "bbb-voting/voting-commons/user-cases"
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
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

	participantRepository := getParticipantRepository(&postgresqlConnector, redisClient, ctx)
	voteRepository := getVotesRepository([]string{os.Getenv("KAFKA_URI")}, redisClient)

	getRoughTotalsUserCase := usercases.NewGetRoughTotalsUserCase(participantRepository, ctx)
	getParticipantsUserCase := usercases.NewGetParticipantsUserCase(participantRepository, ctx)
	castVoteUserCase := usercases.NewCastVoteUserCase(voteRepository, participantRepository, ctx)

	frontendController := controller.NewFrontendController(
		getRoughTotalsUserCase, getParticipantsUserCase, castVoteUserCase,
		ctx, templates, staticFiles,
	)

	serverMux := frontendController.GetServerMux()

	log.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", serverMux); err != nil {
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
func getParticipantRepository(postgresqlConnector *postgresqldatamapper.PostgresqlConnector, redisClient *redis.Client, ctx context.Context) domain.ParticipantRepository {
	var base = postgresqldatamapper.NewParticipantDataMapper(
		postgresqlConnector,
	)

	result, err := redisdatamapper.DecorateParticipantDataMapperWithRedis(base, redisClient, ctx)
	if err != nil {
		log.Fatalf("Faled to load cache: %s", err)
	}

	return result
}
func getVotesRepository(brokers []string, redisClient *redis.Client) domain.VoteRepository {
	var base domain.VoteRepository = kafkadatamapper.NewVoteDataMapper(
		brokers, "votes", 30,
	)
	return redisdatamapper.DecorateVoteDataMapperWithRedis(base, *redisClient)
}
