package usercases

import (
	"bbb-voting/voting-commons/domain"
	"context"
)

type GetRoughTotalsUserCase interface {
	Execute() (map[domain.Participant]int, error)
}
type GetRoughTotalsUserCaseImpl struct {
	voteRepository domain.VoteRepository
	ctx            context.Context
}

func NewGetRoughTotalsUserCaseImpl(voteRepository domain.VoteRepository, ctx context.Context,
) GetRoughTotalsUserCase {
	return GetRoughTotalsUserCaseImpl{voteRepository, ctx}
}

func (userCase GetRoughTotalsUserCaseImpl) Execute() (map[domain.Participant]int, error) {
	return userCase.voteRepository.GetTotalByParticipant(userCase.ctx)
}
