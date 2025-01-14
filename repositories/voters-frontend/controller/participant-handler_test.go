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

var _ = Describe("DashBoardController", func() {
	var (
		ctx        context.Context
		controller FrontendController
	)
	BeforeEach(func() {
		ctx = context.Background()
		controller = NewFrontendController(mocksdatamappers.MockedParticipantDataMapper{}, mocksdatamappers.MockedVotesDataMapper{}, ctx, os.DirFS("../view/templates/"))
	})

	Describe("GetThoroughTotals", func() {
		FIt("Post cast Vote", func() {
			req := httptest.NewRequest("GET", "http://url", nil)

			w := httptest.NewRecorder()

			controller.GetParticipants(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			contentBody := []domain.Participant{}
			json.Unmarshal(body, &contentBody)

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(contentBody).To(Equal(mocksdatamappers.MockedParticipants))
		})
	})
})
