package usercases

import (
	mocksdatamappers "bbb-voting/voting-commons/tests"
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetThoroughTotals", func() {
	var (
		ctx        context.Context
		controller GetThoroughTotalsUserCase
	)
	BeforeEach(func() {
		ctx = context.Background()

		controller = NewGetThoroughTotalsUserCase(mocksdatamappers.MockedParticipantDataMapper{}, ctx)
	})

	It("should return thorough totals", func() {
		result, err := controller.Execute()

		Expect(err).To(BeNil())

		Expect(result.GeneralTotal).To(Equal(len(mocksdatamappers.MockedVotes)))
	})
})
