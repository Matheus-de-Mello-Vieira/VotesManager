package usercases

import (
	"bbb-voting/voting-commons/domain"
	"context"
)

type GetParticipantsUserCase struct {
	participantRepository domain.ParticipantRepository
	ctx                   context.Context
}

func NewGetParticipantsUserCase(participantRepository domain.ParticipantRepository, ctx context.Context,
) GetParticipantsUserCase {
	return GetParticipantsUserCase{participantRepository, ctx}
}

func (userCase GetParticipantsUserCase) Execute() ([]domain.Participant, error) {
	return userCase.participantRepository.FindAll(userCase.ctx)
}
