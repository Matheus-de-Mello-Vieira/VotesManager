package service

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"log"
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
	eventsBatchs, err := register.voteConsumer.GetVoteChan(register.ctx)

	if err != nil {
		log.Fatalf("failed to consume topic: %v", err)
	}

	for votesBulk := range eventsBatchs {
		if len(votesBulk) != 0 {
			log.Printf("Received %d votes", len(votesBulk))
		}
		err = register.voteRepository.SaveMany(*register.ctx, votesBulk)
		if err != nil {
			log.Fatalf("erro on save: %v", err)
		}
	}
}
