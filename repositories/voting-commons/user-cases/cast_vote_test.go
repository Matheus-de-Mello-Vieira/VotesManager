package usercases

import (
	"bbb-voting/voting-commons/domain"
	mocksdatamappers "bbb-voting/voting-commons/tests"
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CastVote", func() {
	var (
		ctx                   context.Context
		sut                   CastVoteUserCase
		participantRepository domain.ParticipantRepository
		voteRepository        domain.VoteRepository
	)
	BeforeEach(func() {
		ctx = context.Background()

		participantRepository = mocksdatamappers.MockedParticipantDataMapper{}
		voteRepository = mocksdatamappers.MockedVotesDataMapper{}

		sut = NewCastVoteUserCaseImpl(voteRepository, participantRepository, ctx)
	})

	It("should cast voter", func() {
		var participantID int = mocksdatamappers.MockedParticipants[0].ParticipantID

		dto := CastVoteDTO{participantID}

		oldLen := len(mocksdatamappers.MockedVotes)
		result, err := sut.Execute(&dto)
		newLen := len(mocksdatamappers.MockedVotes)

		Expect(err).To(BeNil())
		Expect(newLen).To(Equal(oldLen + 1))
		ExpectVoteParticipantEquals(result, participantID)
		ExpectVoteParticipantEquals(&mocksdatamappers.MockedVotes[len(mocksdatamappers.MockedVotes)-1], participantID)
	})

	It("should error when participant don't exists", func() {
		const participantID int = -1

		dto := CastVoteDTO{participantID}

		oldLen := len(mocksdatamappers.MockedVotes)
		result, err := sut.Execute(&dto)
		newLen := len(mocksdatamappers.MockedVotes)

		Expect(err).NotTo(BeNil())
		Expect(errors.Is(err, ErrParticipantNotFound)).To(Equal(true))

		Expect(result).To(BeNil())
		Expect(newLen).To(Equal(oldLen))
	})
})

func ExpectVoteParticipantEquals(vote *domain.Vote, participantID int) {
	Expect(vote.Participant.ParticipantID).To(Equal(participantID))
}
