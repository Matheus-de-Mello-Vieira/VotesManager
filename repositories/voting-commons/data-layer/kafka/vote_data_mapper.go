package kafkadatamapper

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
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
	producer, err := mapper.getProducer()
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
		return err
	}
	defer producer.Close()

	err = mapper.produce(producer, vote)
	if err != nil {
		return err
	}

	log.Printf("1 message sent to Kafka topic %s", mapper.topic)
	return nil
}

func (mapper VoteDataMapper) SaveMany(ctx context.Context, votes []domain.Vote) error {
	producer, err := mapper.getProducer()
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
		return err
	}
	defer producer.Close()

	for _, vote := range votes {
		err = mapper.produce(producer, &vote)
		if err != nil {
			return err
		}
	}

	log.Printf("%d message sent to Kafka topic %s", len(votes), mapper.topic)
	return nil
}

func (mapper VoteDataMapper) getProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	return sarama.NewSyncProducer(mapper.brokers, config)
}

func (mapper VoteDataMapper) produce(producer sarama.SyncProducer, vote *domain.Vote) error {
	voteJSON, err := json.Marshal(vote)
	if err != nil {
		return fmt.Errorf("failed to serialize vote to JSON: %v", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: mapper.topic,
		Key:   sarama.StringEncoder(mapper.getPartitionKey(vote)),
		Value: sarama.StringEncoder(voteJSON),
	}

	_, _, err = producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %v", err)
	}

	return nil
}

func (mapper VoteDataMapper) getPartitionKey(vote *domain.Vote) string {
	return fmt.Sprintf("%d %s", mapper.TruncateUnix(vote), vote.Participant.Name)
}

func (mapper VoteDataMapper) TruncateUnix(vote *domain.Vote) int64 {
	result := vote.Timestamp.Unix() / int64(mapper.aggregateSizeInSeconds)

	return result
}
