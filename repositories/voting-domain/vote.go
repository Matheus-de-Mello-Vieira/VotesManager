package domain

import (
	"time"
)

type Vote struct {
	ID            int
	ParticipantID string
	Timestamp     time.Time
}
