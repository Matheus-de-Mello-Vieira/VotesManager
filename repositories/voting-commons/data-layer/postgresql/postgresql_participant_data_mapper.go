package postgresqldatamapper

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"fmt"
)

type ParticipantRepository struct {
	connector PostgresqlConnector
}

func (mapper *ParticipantRepository) findAll(ctx context.Context) ([]domain.Participant, error) {
	dbpool, err := mapper.connector.openConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer dbpool.Close()

	query := "SELECT participant_id, participant_name FROM participants"

	rows, err := dbpool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query participants: %w", err)
	}
	defer rows.Close()

	var participants []domain.Participant
	for rows.Next() {
		var p domain.Participant
		err := rows.Scan(&p.ParticipantID, &p.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan participant: %w", err)
		}
		participants = append(participants, p)
	}

	return participants, nil
}

// func GetRoughTotals(ctx context.Context) (map[Participant]float64, error) {

// }
// func GetHourlyTotals(ctx context.Context) {

// }
