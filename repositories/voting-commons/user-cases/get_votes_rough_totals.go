package usercases

import (
	"bbb-voting/voting-commons/domain"
	"context"
)

type GetRoughTotalsUserCase struct {
	participantRepository domain.ParticipantRepository
	ctx                   context.Context
}

func NewGetRoughTotalsUserCase(participantRepository domain.ParticipantRepository, ctx context.Context,
) GetRoughTotalsUserCase {
	return GetRoughTotalsUserCase{participantRepository, ctx}
}

func (userCase GetRoughTotalsUserCase) Execute() (map[domain.Participant]int, error) {
	return userCase.participantRepository.GetRoughTotals(userCase.ctx)
}
