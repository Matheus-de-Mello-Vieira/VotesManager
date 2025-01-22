package service

import (
	"bbb-voting/voting-commons/domain"
	mocksdatamappers "bbb-voting/voting-commons/tests"
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetParticipants", func() {
	var (
		ctx                   context.Context
		sut                   GetParticipantsUserCase
		participantRepository domain.ParticipantRepository
	)
	BeforeEach(func() {
		ctx = context.Background()

		participantRepository = mocksdatamappers.MockedParticipantDataMapper{}

		sut = NewGetParticipantsUserCaseImpl(participantRepository, ctx)
	})

	It("should get participants", func() {
		result, err := sut.Execute()

		Expect(err).To(BeNil())
		Expect(result).To(Equal(mocksdatamappers.MockedParticipants))
	})

})
