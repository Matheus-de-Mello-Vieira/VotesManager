package usercases

import (
	"bbb-voting/voting-commons/domain"
	"context"
)

type GetThoroughTotalsUserCase struct {
	participantRepository domain.ParticipantRepository
	ctx                   context.Context
}

func NewGetThoroughTotalsUserCase(participantRepository domain.ParticipantRepository, ctx context.Context,
) GetThoroughTotalsUserCase {
	return GetThoroughTotalsUserCase{participantRepository, ctx}
}

func (userCase GetThoroughTotalsUserCase) Execute() (*domain.ThoroughTotals, error) {
	return userCase.participantRepository.GetThoroughTotals(userCase.ctx)
}
