package domain

import (
	"context"
)

type VoteRepository interface {
	SaveOne(ctx context.Context, vote *Vote) error
	// SaveMany(ctx context.Context, votes []Vote) error
}
