package domain

import (
	"context"
)

type ParticipantRepository interface {
	findAll(ctx context.Context) ([]Participant, error)
	// GetRoughTotals(ctx context.Context) (map[Participant]float64, error)
	// GetHourlyTotals(ctx context.Context)
}
