package main

import (
	kafkaeventconsumer "bbb-voting/votes-register/data_layer"
	"bbb-voting/votes-register/service"
	postgresqldatamapper "bbb-voting/voting-commons/data-layer/postgresql"
	"context"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	timeout, _ := time.ParseDuration("30s")
	batchSize := 1000

	voteConsumer, err := kafkaeventconsumer.NewKafkaVoteConsumer(strings.Split(os.Getenv("KAFKA_URI"), ","), "votes", "events-register", timeout, batchSize)
	if err != nil {
		log.Fatalf("failed to create Kafka consumer: %v", err)
	}

	ctx := context.Background()

	postgresqlConnector := postgresqldatamapper.NewPostgresqlConnector(os.Getenv("POSTGRESQL_URI"))
	voteDataMapper := postgresqldatamapper.NewVoteDataMapper(
		&postgresqlConnector,
	)

	voteRegister := service.NewVoteRegister(voteConsumer, voteDataMapper, &ctx)

	voteRegister.Start()
}
