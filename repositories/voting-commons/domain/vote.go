package domain

import (
	"time"
)

type Vote struct {
	VoteID      int
	participant Participant
	Timestamp   time.Time
}
