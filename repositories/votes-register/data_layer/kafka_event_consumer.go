package kafkaeventconsumer

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

// ConsumerGroupHandler handles messages from the topic
type ConsumerGroupHandler struct {
	eventsBatchs chan<- []domain.Vote
	timeout      time.Duration
	batchSize    int
}

// Setup is run before consuming starts
func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run after consuming ends
func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	close(h.eventsBatchs)
	return nil
}

// ConsumeClaim processes messages from the topic
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var lastMessage *sarama.ConsumerMessage // Keep track of the last message in the batch

	messages := claim.Messages()

	for {
		messageBatch, ended := h.consumeBatch(messages)

		if ended {
			return nil
		}

		if len(messageBatch) > 0 {
			votes, err := h.unmarshalVoteBatch(messageBatch)

			if err != nil {
				return err
			}

			h.eventsBatchs <- votes

			lastMessage = messageBatch[len(messageBatch)-1]
			session.MarkOffset(claim.Topic(), claim.Partition(), lastMessage.Offset+1, "")
		}

	}
}

func (h *ConsumerGroupHandler) consumeBatch(messages <-chan *sarama.ConsumerMessage) ([]*sarama.ConsumerMessage, bool) {
	var result []*sarama.ConsumerMessage
	timer := time.NewTimer(h.timeout)
	defer timer.Stop()

	for i := 0; i < h.batchSize; i++ {
		select {
		case vote, ok := <-messages:
			if !ok {
				// Channel closed
				return result, true
			}
			result = append(result, vote)
		case <-timer.C:
			// Timeout
			return result, false
		}
	}
	return result, false
}
func (h *ConsumerGroupHandler) unmarshalVoteBatch(messages []*sarama.ConsumerMessage) ([]domain.Vote, error) {
	result := make([]domain.Vote, len(messages))

	for i, message := range messages {
		element, err := h.unmarshalVote(message)

		if err != nil {
			return nil, err
		}

		result[i] = *element
	}

	return result, nil
}

func (h *ConsumerGroupHandler) unmarshalVote(message *sarama.ConsumerMessage) (*domain.Vote, error) {
	var vote domain.Vote
	if err := json.Unmarshal(message.Value, &vote); err != nil {
		return nil, err
	}

	return &vote, nil
}

type KafkaVoteConsumer struct {
	brokers   []string
	topic     string
	groupID   string
	timeout   time.Duration
	batchSize int
}

func NewKafkaVoteConsumer(brokers []string, topic string, groupID string, timeout time.Duration, batchSize int) (KafkaVoteConsumer, error) {
	return KafkaVoteConsumer{brokers, topic, groupID, timeout, batchSize}, nil
}

func (kc KafkaVoteConsumer) GetVoteChan(ctx *context.Context) (<-chan []domain.Vote, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Fetch.Max = 10 * 1024 * 1024
	config.Consumer.MaxProcessingTime = kc.timeout

	// Create a new consumer group
	consumerGroup, err := sarama.NewConsumerGroup(kc.brokers, kc.groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	// Create a channel for eventsBatchs
	eventsBatchs := make(chan []domain.Vote)

	// Create a ConsumerGroupHandler and pass the events channel
	handler := &ConsumerGroupHandler{eventsBatchs: eventsBatchs, timeout: kc.timeout, batchSize: kc.batchSize}

	// Start consuming in a separate goroutine
	go func() {
		defer func() {
			if err := consumerGroup.Close(); err != nil {
				log.Printf("failed to close consumer group: %v", err)
			}
			close(eventsBatchs)
		}()

		for {
			if err := consumerGroup.Consume(*ctx, []string{kc.topic}, handler); err != nil {
				log.Printf("Error consuming messages: %v", err)
				break
			}
		}
	}()

	return eventsBatchs, nil
}
