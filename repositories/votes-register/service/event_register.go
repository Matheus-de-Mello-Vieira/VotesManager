package service

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"log"
	"time"
)

type VoteRegister struct {
	voteConsumer   VoteConsumer
	voteRepository domain.VoteRepository
	ctx            *context.Context
}

func NewVoteRegister(voteConsumer VoteConsumer, voteRepository domain.VoteRepository, ctx *context.Context) *VoteRegister {
	return &VoteRegister{
		voteConsumer:   voteConsumer,
		voteRepository: voteRepository,
		ctx:            ctx,
	}
}

func (register *VoteRegister) Start() {
	votes, err := register.voteConsumer.GetVoteChan(register.ctx)

	if err != nil {
		log.Fatalf("failed to consume topic: %v", err)
	}

	timeout, _ := time.ParseDuration("30s")
	for {
		votesBulk := consumeWithTimeout(votes, 1000, timeout)

		if len(votesBulk) != 0 {
			log.Printf("Received %d votes", len(votesBulk))
		}
		err = register.voteRepository.SaveMany(*register.ctx, votesBulk)
		if err != nil {
			log.Println("erro on save: %w")
		}
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
