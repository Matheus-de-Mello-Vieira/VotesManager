package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"bbb-voting/voting-commons/domain"
	mocksdatamappers "bbb-voting/voting-commons/tests"
	"context"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"
)

var _ = Describe("VotesController", func() {
	var (
		ctx                   context.Context
		controller            FrontendController
		participantRepository domain.ParticipantRepository
	)
	BeforeEach(func() {
		ctx = context.Background()
		participantRepository = mocksdatamappers.MockedParticipantDataMapper{}
		controller = NewFrontendController(participantRepository, mocksdatamappers.MockedVotesDataMapper{}, ctx, os.DirFS("../view/templates/"))
	})

	Describe("GetVotesRoughTotalsHandler", func() {
		It("MustShowTheData", func() {
			req := httptest.NewRequest("GET", "http://example.com/votes/totals/rough", nil)

			w := httptest.NewRecorder()

			controller.GetVotesRoughTotalsHandler(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			actualResult := map[string]int{}
			json.Unmarshal(body, &actualResult)

			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			expectedResult, err := participantRepository.GetRoughTotals(ctx)
			if err != nil {
				Fail(fmt.Sprint("error on GetRoughTotals:", err))
			}

			Expect(len(actualResult)).To(Equal(len(expectedResult)))

			for participant, total := range expectedResult {
				Expect(actualResult[participant.Name]).To(Equal(total))
			}
		})
	})

	Describe("CastVote", func() {
		It("vote", func() {
			const participantID int = 1

			inputData := map[string]int{"participant_id": participantID}
			inputBody, _ := json.Marshal(inputData)
			req := httptest.NewRequest("POST", "http://example.com/votes", bytes.NewBuffer(inputBody))
			req.Header.Add("Content-Type", "application/json")

			w := httptest.NewRecorder()

			oldLen := len(mocksdatamappers.MockedVotes)
			controller.PostVoteHandler(w, req)

			newLen := len(mocksdatamappers.MockedVotes)

			resp := w.Result()
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))
			Expect(newLen).To(Equal(oldLen + 1))
			Expect(mocksdatamappers.MockedVotes[len(mocksdatamappers.MockedVotes)-1].Participant.ParticipantID).To(Equal(participantID))
		})
	})
})
