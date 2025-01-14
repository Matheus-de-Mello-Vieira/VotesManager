package kafkaeventconsumer

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

// ConsumerGroupHandler handles messages from the topic
type ConsumerGroupHandler struct {
	events chan<- domain.Vote
}

// Setup is run before consuming starts
func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run after consuming ends
func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	close(h.events)
	return nil
}

// ConsumeClaim processes messages from the topic
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var vote domain.Vote
		if err := json.Unmarshal(message.Value, &vote); err != nil {
			log.Printf("failed to unmarshal message: %v", err)
			continue
		}
		h.events <- vote                 // Send the event to the channel
		session.MarkMessage(message, "") // Mark message as processed
	}
	return nil
}

type KafkaVoteConsumer struct {
	brokers []string
	topic   string
	groupID string
}

func NewKafkaVoteConsumer(brokers []string, topic string, groupID string) (KafkaVoteConsumer, error) {
	return KafkaVoteConsumer{brokers, topic, groupID}, nil
}

func (kc KafkaVoteConsumer) GetVoteChan(ctx *context.Context) (<-chan domain.Vote, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest // Start from the latest message

	// Create a new consumer group
	consumerGroup, err := sarama.NewConsumerGroup(kc.brokers, kc.groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	// Create a channel for events
	events := make(chan domain.Vote)

	// Create a ConsumerGroupHandler and pass the events channel
	handler := &ConsumerGroupHandler{events: events}

	// Start consuming in a separate goroutine
	go func() {
		defer func() {
			if err := consumerGroup.Close(); err != nil {
				log.Printf("failed to close consumer group: %v", err)
			}
			close(events)
		}()

		for {
			if err := consumerGroup.Consume(*ctx, []string{kc.topic}, handler); err != nil {
				log.Printf("Error consuming messages: %v", err)
				break
			}
		}
	}()

	return events, nil
}
