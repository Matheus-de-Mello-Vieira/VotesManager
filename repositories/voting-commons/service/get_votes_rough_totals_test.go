package service

import (
	"bbb-voting/voting-commons/domain"
	mocksdatamappers "bbb-voting/voting-commons/tests"
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetRoughTotals", func() {
	var (
		ctx            context.Context
		sut            GetRoughTotalsUserCase
		voteRepository domain.VoteRepository
	)
	BeforeEach(func() {
		ctx = context.Background()

		voteRepository = mocksdatamappers.MockedVotesDataMapper{}

		sut = NewGetRoughTotalsUserCaseImpl(voteRepository, ctx)
	})

	It("should return thorough totals", func() {
		actualResult, err := sut.Execute()
		Expect(err).To(BeNil())

		expectedResult, err := voteRepository.GetTotalByParticipant(ctx)
		if err != nil {
			Fail(fmt.Sprint("error on GetTotalByParticipant:", err))
		}

		Expect(len(actualResult)).To(Equal(len(expectedResult)))
		for participant, total := range expectedResult {
			Expect(actualResult[participant]).To(Equal(total))
		}
	})
})
