package usercases

import (
	"bbb-voting/voting-commons/domain"
	mocksdatamappers "bbb-voting/voting-commons/tests"
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetThoroughTotals", func() {
	var (
		ctx                   context.Context
		sut                   GetThoroughTotalsUserCase
		participantRepository domain.ParticipantRepository
	)
	BeforeEach(func() {
		ctx = context.Background()

		participantRepository = mocksdatamappers.MockedParticipantDataMapper{}

		sut = NewGetThoroughTotalsUserCaseImpl(participantRepository, ctx)
	})

	It("should return thorough totals", func() {
		result, err := sut.Execute()

		Expect(err).To(BeNil())

		Expect(result.GeneralTotal).To(Equal(len(mocksdatamappers.MockedVotes)))
	})
})
