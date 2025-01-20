package postgresqldatamapper

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"fmt"
	"strings"
)

type ParticipantDataMapper struct {
	connector *PostgresqlConnector
}

func NewParticipantDataMapper(connector *PostgresqlConnector) ParticipantDataMapper {
	return ParticipantDataMapper{connector}
}

func (mapper ParticipantDataMapper) FindAll(ctx context.Context) ([]domain.Participant, error) {
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
		p.Name = strings.TrimSpace(p.Name)
		participants = append(participants, p)
	}

	return participants, nil
}

func (mapper ParticipantDataMapper) FindByID(ctx context.Context, id int) (*domain.Participant, error) {
	dbpool, err := mapper.connector.openConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer dbpool.Close()

	query := "SELECT participant_id, participant_name FROM participants WHERE participant_id = $1"

	var p domain.Participant
	err = dbpool.QueryRow(ctx, query, id).Scan(&p.ParticipantID, &p.Name)
	p.Name = strings.TrimSpace(p.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get participant: %w", err)
	}

	return &p, nil
}
