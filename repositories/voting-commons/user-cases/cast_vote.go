package usercases

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"errors"
	"time"
)

type CastVoteUserCase struct {
	voteRepository        domain.VoteRepository
	participantRepository domain.ParticipantRepository
	ctx                   context.Context
}

type CastVoteDTO struct {
	ParticipantID int `json:"participant_id"`
}

var ErrParticipantNotFound = errors.New("participant not found")

func NewCastVoteUserCase(voteRepository domain.VoteRepository, participantRepository domain.ParticipantRepository, ctx context.Context,
) CastVoteUserCase {
	return CastVoteUserCase{voteRepository, participantRepository, ctx}
}

func (userCase CastVoteUserCase) Execute(castVoteDto *CastVoteDTO) (*domain.Vote, error) {
	participant, _ := userCase.participantRepository.FindByID(userCase.ctx, castVoteDto.ParticipantID)

	if participant == nil {
		return nil, ErrParticipantNotFound
	}

	vote := domain.Vote{Participant: *participant, Timestamp: time.Now()}

	userCase.voteRepository.SaveOne(userCase.ctx, &vote)

	return &vote, nil
}
