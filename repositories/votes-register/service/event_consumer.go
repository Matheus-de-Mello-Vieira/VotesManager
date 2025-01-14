package service

import (
	"bbb-voting/voting-commons/domain"
	"context"
)

type VoteConsumer interface {
	GetVoteChan(ctx *context.Context) (<-chan domain.Vote, error)
}
