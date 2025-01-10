package mocksdatamappers

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"fmt"
)

var IsaacNewton = domain.Participant{ParticipantID: 1,
	Name: "Isaac Newton"}

var AlbertEinstein = domain.Participant{ParticipantID: 2,
	Name: "Albert Einstein"}
var MarieCurie = domain.Participant{
	ParticipantID: 3,
	Name:          "Marie Curie",
}

var MockedParticipants = []domain.Participant{
	IsaacNewton, AlbertEinstein, MarieCurie,
}

type MockedParticipantDataMapper struct{}

func (mapper MockedParticipantDataMapper) FindAll(ctx context.Context) ([]domain.Participant, error) {
	return MockedParticipants, nil
}

func (mapper MockedParticipantDataMapper) FindByID(ctx context.Context, id int) (*domain.Participant, error) {
	for _, participant := range MockedParticipants {
		if participant.ParticipantID == id {
			return &participant, nil
		}
	}

	return nil, fmt.Errorf("failed to get participant with id %v", id)
}

func (mapper MockedParticipantDataMapper) GetRoughTotals(ctx context.Context) (map[domain.Participant]int, error) {
	return getVotesByParticipant()
}

func (mapper MockedParticipantDataMapper) GetThoroughTotals(ctx context.Context) (*domain.ThoroughTotals, error) {
	generalTotal, err1 := getGeneralTotal()
	if err1 != nil {
		return nil, err1
	}

	totalByParticipant, err2 := getVotesByParticipant()
	if err2 != nil {
		return nil, err2
	}

	totalByHour, err3 := getVotesByHour()
	if err3 != nil {
		return nil, err3
	}

	result := domain.ThoroughTotals{GeneralTotal: *generalTotal, TotalByHour: totalByHour, TotalByParticipant: totalByParticipant}

	return &result, nil
}

func getVotesByParticipant() (map[domain.Participant]int, error) {
	var result = map[domain.Participant]int{}

	for _, vote := range MockedVotes {
		if _, exists := result[vote.Participant]; exists {
			result[vote.Participant] = 0
		}
		result[vote.Participant]++
	}

	return result, nil
}
