package kafkadatamapper

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
)

type VoteDataMapper struct {
	brokers                []string
	topic                  string
	aggregateSizeInSeconds int
}

func NewVoteDataMapper(brokers []string, topic string, aggregateSizeInSeconds int) VoteDataMapper {
	return VoteDataMapper{brokers, topic, aggregateSizeInSeconds}
}

func (mapper VoteDataMapper) SaveOne(ctx context.Context, vote *domain.Vote) error {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(mapper.brokers, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// Serialize the vote to JSON
	voteJSON, err := json.Marshal(vote)
	if err != nil {
		log.Fatalf("Failed to serialize vote to JSON: %v", err)
	}

	// Publish the vote JSON to the Kafka topic
	msg := &sarama.ProducerMessage{
		Topic: mapper.topic,
		Key:   sarama.StringEncoder(mapper.getPartitionKey(vote)),
		Value: sarama.StringEncoder(voteJSON),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatalf("Failed to send message to Kafka: %v", err)
	}

	log.Printf("Message sent to Kafka topic %s (partition: %d, offset: %d)", mapper.topic, partition, offset)
	return nil
}

func (mapper VoteDataMapper) getPartitionKey(vote *domain.Vote) string {
	return fmt.Sprintf("%d %s", mapper.TruncateUnix(vote), vote.Participant.Name)
}

func (mapper VoteDataMapper) TruncateUnix(vote *domain.Vote) int64 {
	result := vote.Timestamp.Unix() / int64(mapper.aggregateSizeInSeconds)

	return result
}
