package domain

import (
	"context"
)

type VoteRepository interface {
	SaveOne(ctx context.Context, vote *Vote) error
	SaveMany(ctx context.Context, votes []Vote) error

	GetGeneralTotal(ctx context.Context) (int, error)
	GetTotalByHour(ctx context.Context) ([]TotalByHour, error)
	GetTotalByParticipant(ctx context.Context) (map[Participant]int, error)
}
