package domain

import (
	"time"
)

type Vote struct {
	VoteID      int
	Participant Participant
	Timestamp   time.Time
}
