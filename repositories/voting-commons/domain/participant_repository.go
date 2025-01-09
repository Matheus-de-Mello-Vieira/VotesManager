package domain

type ParticipantRepository interface {
	GetRoughTotals() (map[Participant]float64, error)
	GetHourlyTotals()
}
