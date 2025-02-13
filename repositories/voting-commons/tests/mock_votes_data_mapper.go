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

func (mapper MockedVotesDataMapper) SaveMany(ctx context.Context, votes []domain.Vote) error {
	for _, vote := range votes {
		mapper.SaveOne(ctx, &vote)
	}
	return nil
}

func (mapper MockedVotesDataMapper) GetGeneralTotal(ctx context.Context) (int, error) {
	return len(MockedVotes), nil
}

func (mapper MockedVotesDataMapper) GetTotalByHour(ctx context.Context) ([]domain.TotalByHour, error) {
	var result []domain.TotalByHour

	for hour, total := range getMapByHour() {
		result = append(result, domain.TotalByHour{
			Total: total,
			Hour:  hour,
		})
	}

	return result, nil
}

func getMapByHour() map[time.Time]int {
	var result = map[time.Time]int{}

	for _, vote := range MockedVotes {
		if _, exists := result[vote.GetHour()]; exists {
			result[vote.GetHour()] = 0
		}
		result[vote.GetHour()]++
	}

	return result
}

func (mapper MockedVotesDataMapper) GetTotalByParticipant(ctx context.Context) (map[domain.Participant]int, error) {
	var result = map[domain.Participant]int{}

	for _, vote := range MockedVotes {
		if _, exists := result[vote.Participant]; exists {
			result[vote.Participant] = 0
		}
		result[vote.Participant]++
	}

	return result, nil
}
