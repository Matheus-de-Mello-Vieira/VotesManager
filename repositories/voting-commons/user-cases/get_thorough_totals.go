package usercases

import (
	"bbb-voting/voting-commons/domain"
	"context"
)
type GetThoroughTotalsUserCase interface {
	Execute() (*domain.ThoroughTotals, error)
}
type GetThoroughTotalsUserCaseImpl struct {
	participantRepository domain.ParticipantRepository
	ctx                   context.Context
}

func NewGetThoroughTotalsUserCaseImpl(participantRepository domain.ParticipantRepository, ctx context.Context,
) GetThoroughTotalsUserCaseImpl {
	return GetThoroughTotalsUserCaseImpl{participantRepository, ctx}
}

func (userCase GetThoroughTotalsUserCaseImpl) Execute() (*domain.ThoroughTotals, error) {
	return userCase.participantRepository.GetThoroughTotals(userCase.ctx)
}
