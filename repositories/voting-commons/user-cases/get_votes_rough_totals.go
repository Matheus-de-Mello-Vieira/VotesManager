package usercases

import (
	"bbb-voting/voting-commons/domain"
	"context"
)

type GetRoughTotalsUserCase interface {
	Execute() (map[domain.Participant]int, error)
}
type GetRoughTotalsUserCaseImpl struct {
	participantRepository domain.ParticipantRepository
	ctx                   context.Context
}

func NewGetRoughTotalsUserCaseImpl(participantRepository domain.ParticipantRepository, ctx context.Context,
) GetRoughTotalsUserCase {
	return GetRoughTotalsUserCaseImpl{participantRepository, ctx}
}

func (userCase GetRoughTotalsUserCaseImpl) Execute() (map[domain.Participant]int, error) {
	return userCase.participantRepository.GetRoughTotals(userCase.ctx)
}
