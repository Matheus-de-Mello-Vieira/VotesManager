package controller

import (
	"encoding/json"
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

var _ = Describe("ParticipantsController", func() {
	var (
		ctx        context.Context
		controller FrontendController
	)
	BeforeEach(func() {
		ctx = context.Background()
		controller = NewFrontendController(mocksdatamappers.MockedParticipantDataMapper{}, mocksdatamappers.MockedVotesDataMapper{}, ctx, os.DirFS("../view/templates/"), os.DirFS("../view/static/"))
	})

	Describe("GetParticipantsHandler", func() {
		It("MustShowTheList", func() {
			req := httptest.NewRequest("GET", "http://example.com/participants", nil)

			w := httptest.NewRecorder()

			controller.GetParticipantsHandler(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			contentBody := []domain.Participant{}
			json.Unmarshal(body, &contentBody)

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(contentBody).To(Equal(mocksdatamappers.MockedParticipants))
		})
	})
})
