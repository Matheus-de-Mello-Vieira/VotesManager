package domain

import (
	"context"
)

type TotalByHour struct {
	Total int `json:"total"`
	Hour  int `json:"hour"`
}
type ThoroughTotals struct {
	GeneralTotal       int
	TotalByHour        []TotalByHour
	TotalByParticipant map[Participant]int
}

type ParticipantRepository interface {
	FindAll(ctx context.Context) ([]Participant, error)
	FindByID(ctx context.Context, id int) (*Participant, error)
	GetRoughTotals(ctx context.Context) (map[Participant]int, error)
	GetThoroughTotals(ctx context.Context) (*ThoroughTotals, error)
}
