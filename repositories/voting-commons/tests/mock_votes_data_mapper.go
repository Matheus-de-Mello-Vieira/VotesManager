package mocksdatamappers

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"time"
)

func newDateByHourMinute(hour int, minute int) time.Time {
	return time.Date(2025, 1, 1, hour, minute, 1, 1, time.UTC)
}

var MockedVotes = []domain.Vote{
	{
		VoteID:      1,
		Participant: IsaacNewton,
		Timestamp:   newDateByHourMinute(1, 0),
	},
	{
		VoteID:      2,
		Participant: AlbertEinstein,
		Timestamp:   newDateByHourMinute(1, 0),
	},
	{
		VoteID:      3,
		Participant: MarieCurie,
		Timestamp:   newDateByHourMinute(1, 0),
	},
	{
		VoteID:      4,
		Participant: MarieCurie,
		Timestamp:   newDateByHourMinute(1, 0),
	},
	{
		VoteID:      5,
		Participant: MarieCurie,
		Timestamp:   newDateByHourMinute(1, 0),
	},
}

type MockedVotesDataMapper struct{}

func (mapper MockedVotesDataMapper) SaveOne(ctx context.Context, vote *domain.Vote) error {
	MockedVotes = append(MockedVotes, *vote)
	vote.VoteID = MockedVotes[len(MockedVotes)-1].VoteID + 1
	return nil
}

func getGeneralTotal() (*int, error) {
	result := len(MockedVotes)
	return &result, nil
}

func getVotesByHour() ([]domain.TotalByHour, error) {
	var result []domain.TotalByHour

	for hour, total := range getMapByHour() {
		result = append(result, domain.TotalByHour{
			Total: total,
			Hour:  hour,
		})
	}

	return result, nil
}

func getMapByHour() map[int]int {
	var result = map[int]int{}

	for _, vote := range MockedVotes {
		if _, exists := result[vote.Timestamp.Hour()]; exists {
			result[vote.Timestamp.Hour()] = 0
		}
		result[vote.Timestamp.Hour()]++
	}

	return result
}
