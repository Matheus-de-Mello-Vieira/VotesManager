package controller

import (
	mocksdatamappers "bbb-voting/voting-commons/tests"
	"context"
	"encoding/json"
	"io"

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
		controller = NewFrontendController(mocksdatamappers.MockedParticipantDataMapper{}, ctx)
	})

	Describe("GetThoroughTotals", func() {
		FIt("should return thorough totals", func() {
			req := httptest.NewRequest("GET", "http://example.com/votes/thorough", nil)
			w := httptest.NewRecorder()

			controller.GetThoroughTotals(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			Expect(resp.StatusCode).To(Equal(200))

			contentBody := ThoroughTotalsResponseModel{}
			json.Unmarshal(body, &contentBody)

			Expect(contentBody.GeneralTotal).To(Equal(len(mocksdatamappers.MockedVotes)))
		})
	})
})
