package domain

import (
	"time"
)

const SECONDS_IN_HOUR = int64(60 * 60)

type Vote struct {
	VoteID      int
	Participant Participant
	Timestamp   time.Time
}

func (vote *Vote) GetHour() time.Time {
	resultTimestamp := (vote.Timestamp.Unix() / SECONDS_IN_HOUR) * SECONDS_IN_HOUR

	return time.Unix(resultTimestamp, 0)
}
