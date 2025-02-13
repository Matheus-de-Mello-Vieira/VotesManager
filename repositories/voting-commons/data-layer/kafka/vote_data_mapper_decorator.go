package kafkadatamapper

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
)

type VoteDataMapperKafkaDecorator struct {
	brokers                []string
	topic                  string
	base                   domain.VoteRepository
}

func DecorateVoteDataRepository(base domain.VoteRepository, brokers []string, topic string) VoteDataMapperKafkaDecorator {
	return VoteDataMapperKafkaDecorator{brokers, topic, base}
}

func (mapper VoteDataMapperKafkaDecorator) SaveOne(ctx context.Context, vote *domain.Vote) error {
	producer, err := mapper.getProducer()
	if err != nil {
		return fmt.Errorf("failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	err = mapper.produce(producer, vote)
	if err != nil {
		return err
	}

	return nil
}

func (mapper VoteDataMapperKafkaDecorator) SaveMany(ctx context.Context, votes []domain.Vote) error {
	producer, err := mapper.getProducer()
	if err != nil {
		return fmt.Errorf("failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	for _, vote := range votes {
		err = mapper.produce(producer, &vote)
		if err != nil {
			return err
		}
	}

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
		Key:   nil,
		Value: sarama.StringEncoder(voteJSON),
	}

	_, _, err = producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %v", err)
	}

	return nil
}

func (mapper VoteDataMapperKafkaDecorator) GetGeneralTotal(ctx context.Context) (int, error) {
	return mapper.base.GetGeneralTotal(ctx)
}
func (mapper VoteDataMapperKafkaDecorator) GetTotalByHour(ctx context.Context) ([]domain.TotalByHour, error) {
	return mapper.base.GetTotalByHour(ctx)
}
func (mapper VoteDataMapperKafkaDecorator) GetTotalByParticipant(ctx context.Context) (map[domain.Participant]int, error) {
	return mapper.base.GetTotalByParticipant(ctx)
}
