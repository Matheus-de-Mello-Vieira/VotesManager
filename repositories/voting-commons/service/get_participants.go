package service

import (
	"bbb-voting/voting-commons/domain"
	"context"
)

type GetParticipantsUserCase interface {
	Execute() ([]domain.Participant, error)
}
type GetParticipantsUserCaseImpl struct {
	participantRepository domain.ParticipantRepository
	ctx                   context.Context
}

func NewGetParticipantsUserCaseImpl(participantRepository domain.ParticipantRepository, ctx context.Context,
) GetParticipantsUserCaseImpl {
	return GetParticipantsUserCaseImpl{participantRepository, ctx}
}

func (userCase GetParticipantsUserCaseImpl) Execute() ([]domain.Participant, error) {
	return userCase.participantRepository.FindAll(userCase.ctx)
}
