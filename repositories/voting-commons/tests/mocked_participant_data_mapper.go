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
