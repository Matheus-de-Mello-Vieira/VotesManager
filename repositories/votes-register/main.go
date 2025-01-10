package main

import (
	kafkaeventconsumer "bbb-voting/votes-register/data_layer"
	postgresqldatamapper "bbb-voting/voting-commons/data-layer/postgresql"
	"bbb-voting/voting-commons/domain"
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	consumer, err := kafkaeventconsumer.NewKafkaVoteConsumer([]string{os.Getenv("KAFKA_URI")}, "events", "events-register")
	if err != nil {
		log.Fatalf("failed to create Kafka consumer: %v", err)
	}

	ctx := context.Background()
	votes, error := consumer.Consume(ctx)

	postgresqlConnector := postgresqldatamapper.NewPostgresqlConnector(os.Getenv("POSTGRESQL_URI"))
	voteDataMapper := postgresqldatamapper.NewVoteDataMapper(
		postgresqlConnector,
	)

	if error != nil {
		log.Fatalf("failed to consume topic: %v", err)
	}

	timeout, _ := time.ParseDuration("30s")
	for {
		votesBulk := consumeWithTimeout(votes, 1000, timeout)
		fmt.Printf("Received %d votes", len(votesBulk))
		voteDataMapper.SaveMany(ctx, votesBulk)
	}
}

func consumeWithTimeout(ch <-chan domain.Vote, maxItems int, timeout time.Duration) []domain.Vote {
	var result []domain.Vote
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for i := 0; i < maxItems; i++ {
		select {
		case vote, ok := <-ch:
			if !ok {
				// Channel closed
				return result
			}
			result = append(result, vote)
		case <-timer.C:
			// Timeout
			return result
		}
	}
	return result
}
