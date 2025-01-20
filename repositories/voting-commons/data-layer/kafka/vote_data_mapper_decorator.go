package kafkadatamapper

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type VoteDataMapperKafkaDecorator struct {
	brokers                []string
	topic                  string
	aggregateSizeInSeconds int
	base                   domain.VoteRepository
}

func DecorateVoteDataRepository(base domain.VoteRepository, brokers []string, topic string, aggregateSizeInSeconds int) VoteDataMapperKafkaDecorator {
	return VoteDataMapperKafkaDecorator{brokers, topic, aggregateSizeInSeconds, base}
}

func (mapper VoteDataMapperKafkaDecorator) SaveOne(ctx context.Context, vote *domain.Vote) error {
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

func (mapper VoteDataMapperKafkaDecorator) SaveMany(ctx context.Context, votes []domain.Vote) error {
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

func (mapper VoteDataMapperKafkaDecorator) getProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	return sarama.NewSyncProducer(mapper.brokers, config)
}

func (mapper VoteDataMapperKafkaDecorator) produce(producer sarama.SyncProducer, vote *domain.Vote) error {
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

func (mapper VoteDataMapperKafkaDecorator) getPartitionKey(vote *domain.Vote) string {
	return fmt.Sprintf("%d %s", mapper.TruncateUnix(vote), vote.Participant.Name)
}

func (mapper VoteDataMapperKafkaDecorator) TruncateUnix(vote *domain.Vote) int64 {
	result := vote.Timestamp.Unix() / int64(mapper.aggregateSizeInSeconds)

	return result
}

func (mapper VoteDataMapperKafkaDecorator)  GetGeneralTotal(ctx context.Context) (int, error) {
	return mapper.base.GetGeneralTotal(ctx)
}
func (mapper VoteDataMapperKafkaDecorator) GetTotalByHour(ctx context.Context) ([]domain.TotalByHour, error) {
	return mapper.base.GetTotalByHour(ctx)
}
func (mapper VoteDataMapperKafkaDecorator) GetTotalByParticipant(ctx context.Context) (map[domain.Participant]int, error) {
	return mapper.base.GetTotalByParticipant(ctx)
}