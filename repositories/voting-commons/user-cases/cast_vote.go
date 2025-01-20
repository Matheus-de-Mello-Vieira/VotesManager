package usercases

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"errors"
	"time"
)

type CastVoteUserCase interface {
	Execute(castVoteDto *CastVoteDTO) (*domain.Vote, error)
}

type CastVoteUserCaseImpl struct {
	voteRepository        domain.VoteRepository
	participantRepository domain.ParticipantRepository
	ctx                   context.Context
}

type CastVoteDTO struct {
	ParticipantID int `json:"participant_id"`
}

var ErrParticipantNotFound = errors.New("participant not found")

func NewCastVoteUserCaseImpl(voteRepository domain.VoteRepository, participantRepository domain.ParticipantRepository, ctx context.Context,
) CastVoteUserCaseImpl {
	return CastVoteUserCaseImpl{voteRepository, participantRepository, ctx}
}

func (userCase CastVoteUserCaseImpl) Execute(castVoteDto *CastVoteDTO) (*domain.Vote, error) {
	participant, _ := userCase.participantRepository.FindByID(userCase.ctx, castVoteDto.ParticipantID)

	if participant == nil {
		return nil, ErrParticipantNotFound
	}

	vote := domain.Vote{Participant: *participant, Timestamp: time.Now()}

	userCase.voteRepository.SaveOne(userCase.ctx, &vote)

	return &vote, nil
}
