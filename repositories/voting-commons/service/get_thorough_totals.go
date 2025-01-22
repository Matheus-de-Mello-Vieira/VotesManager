package service

import (
	"bbb-voting/voting-commons/domain"
	"context"
)

type GetThoroughTotalsUserCase interface {
	Execute() (*domain.ThoroughTotals, error)
}
type GetThoroughTotalsUserCaseImpl struct {
	voteRepository domain.VoteRepository
	ctx            context.Context
}

func NewGetThoroughTotalsUserCaseImpl(voteRepository domain.VoteRepository, ctx context.Context,
) GetThoroughTotalsUserCaseImpl {
	return GetThoroughTotalsUserCaseImpl{voteRepository, ctx}
}

func (userCase GetThoroughTotalsUserCaseImpl) Execute() (*domain.ThoroughTotals, error) {
	generalTotal, err := userCase.voteRepository.GetGeneralTotal(userCase.ctx)
	if err != nil {
		return nil, err
	}

	totalByParticipant, err := userCase.voteRepository.GetTotalByParticipant(userCase.ctx)
	if err != nil {
		return nil, err
	}

	totalByHour, err := userCase.voteRepository.GetTotalByHour(userCase.ctx)
	if err != nil {
		return nil, err
	}

	result := domain.ThoroughTotals{GeneralTotal: generalTotal, TotalByHour: totalByHour, TotalByParticipant: totalByParticipant}

	return &result, nil
}
