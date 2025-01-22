package main

import (
	"bbb-voting/voters-frontend/controller"
	kafkadatamapper "bbb-voting/voting-commons/data-layer/kafka"
	localdatamapper "bbb-voting/voting-commons/data-layer/local-cache"
	postgresqldatamapper "bbb-voting/voting-commons/data-layer/postgresql"
	redisdatamapper "bbb-voting/voting-commons/data-layer/redis"
	"bbb-voting/voting-commons/domain"
	"bbb-voting/voting-commons/service"
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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

	participantRepository := getParticipantRepository(&postgresqlConnector, redisClient)
	voteRepository := getVotesRepository(&postgresqlConnector, participantRepository, strings.Split(os.Getenv("KAFKA_URI"), ","), redisClient)

	getRoughTotalsUserCase := service.NewGetRoughTotalsUserCaseImpl(voteRepository, ctx)
	getParticipantsUserCase := service.NewGetParticipantsUserCaseImpl(participantRepository, ctx)
	castVoteUserCase := service.NewCastVoteUserCaseImpl(voteRepository, participantRepository, ctx)

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
func getParticipantRepository(postgresqlConnector *postgresqldatamapper.PostgresqlConnector, redisClient *redis.Client) domain.ParticipantRepository {
	var result domain.ParticipantRepository

	redisCacheTTL, _ := time.ParseDuration("12h")
	localCacheTTL, _ := time.ParseDuration("6h")

	result = postgresqldatamapper.NewParticipantDataMapper(postgresqlConnector)
	result = redisdatamapper.DecorateParticipantRepository(result, redisClient, redisCacheTTL)
	result = localdatamapper.DecorateParticipantRepository(result, localCacheTTL)

	return result
}
func getVotesRepository(postgresqlConnector *postgresqldatamapper.PostgresqlConnector, participantRepository domain.ParticipantRepository, brokers []string, redisClient *redis.Client) domain.VoteRepository {
	var result domain.VoteRepository

	result = postgresqldatamapper.NewVoteDataMapper(postgresqlConnector)
	result = kafkadatamapper.DecorateVoteDataRepository(result, brokers, "votes")
	result = redisdatamapper.DecorateVoteDataRepository(result, *redisClient, participantRepository)

	return result
}
