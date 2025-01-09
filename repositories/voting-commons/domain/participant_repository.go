package domain

import (
	"context"
)

type ParticipantRepository interface {
	FindAll(ctx context.Context) ([]Participant, error)
	FindByID(ctx context.Context, id int) (*Participant, error)
	GetRoughTotals(ctx context.Context) (map[Participant]int, error)
	// GetHourlyTotals(ctx context.Context)
}
