package main

import (
	kafkaeventconsumer "bbb-voting/votes-register/data_layer"
	"bbb-voting/votes-register/service"
	postgresqldatamapper "bbb-voting/voting-commons/data-layer/postgresql"
	"context"
	"log"
	"os"
)

func main() {
	voteConsumer, err := kafkaeventconsumer.NewKafkaVoteConsumer([]string{os.Getenv("KAFKA_URI")}, "votes", "events-register")
	if err != nil {
		log.Fatalf("failed to create Kafka consumer: %v", err)
	}

	ctx := context.Background()

	postgresqlConnector := postgresqldatamapper.NewPostgresqlConnector(os.Getenv("POSTGRESQL_URI"))
	voteDataMapper := postgresqldatamapper.NewVoteDataMapper(
		postgresqlConnector,
	)

	voteRegister := service.NewVoteRegister(voteConsumer, voteDataMapper, &ctx)

	voteRegister.Start()
}
