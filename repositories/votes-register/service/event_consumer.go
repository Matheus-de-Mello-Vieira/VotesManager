package service

import (
	"bbb-voting/voting-commons/domain"
	"context"
)

type VoteConsumer interface {
	Consume(ctx context.Context) ([]domain.Vote, error)
}

