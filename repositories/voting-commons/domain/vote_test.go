package domain

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vote", func() {
	Describe("GetHour", func() {
		testGetHour(time.Date(2022, 1, 1, 2, 4, 0, 0, time.Local), time.Date(2022, 1, 1, 2, 0, 0, 0, time.Local))
		testGetHour(time.Date(2023, 1, 1, 3, 0, 0, 0, time.UTC), time.Date(2023, 1, 1, 3, 0, 0, 0, time.UTC))
	})

})

const TimeStampFormat = "15:04:05 MST"

func testGetHour(voteTimeStamp time.Time, expectedResult time.Time) {
	testFunc := func() {
		sut := Vote{
			VoteID: 1,
			Participant: Participant{ParticipantID: 1,
				Name: "Isaac Newton"},
			Timestamp: voteTimeStamp,
		}

		actualResult := sut.GetHour()

		Expect(actualResult).To(Equal(expectedResult))
	}

	text := fmt.Sprintf("%s -> %s", voteTimeStamp.Format(TimeStampFormat), expectedResult.Format(TimeStampFormat))

	It(text, testFunc)
}
