package service

import (
	"bbb-voting/voting-commons/domain"
	mocksdatamappers "bbb-voting/voting-commons/tests"
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetThoroughTotals", func() {
	var (
		ctx            context.Context
		sut            GetThoroughTotalsUserCase
		voteRepository domain.VoteRepository
	)
	BeforeEach(func() {
		ctx = context.Background()

		voteRepository = mocksdatamappers.MockedVotesDataMapper{}

		sut = NewGetThoroughTotalsUserCaseImpl(voteRepository, ctx)
	})

	It("should return thorough totals", func() {
		result, err := sut.Execute()

		Expect(err).To(BeNil())

		Expect(result.GeneralTotal).To(Equal(len(mocksdatamappers.MockedVotes)))
	})
})
