package postgresqldatamapper

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
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

func (mapper ParticipantDataMapper) GetRoughTotals(ctx context.Context) (map[domain.Participant]int, error) {
	dbpool, err := mapper.connector.openConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer dbpool.Close()

	query := `select * from rough_totals`
	return getManyFromQuery(dbpool, ctx, query)
}

func (mapper ParticipantDataMapper) GetThoroughTotals(ctx context.Context) (*domain.ThoroughTotals, error) {
	dbpool, err := mapper.connector.openConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer dbpool.Close()

	generalTotal, err := getGeneralTotal(dbpool, ctx)
	if err != nil {
		return nil, err
	}

	totalByParticipant, err := getVotesByParticipant(dbpool, ctx)
	if err != nil {
		return nil, err
	}

	totalByHour, err := getVotesByHour(dbpool, ctx)
	if err != nil {
		return nil, err
	}

	result := domain.ThoroughTotals{GeneralTotal: *generalTotal, TotalByHour: totalByHour, TotalByParticipant: totalByParticipant}

	return &result, nil
}

func getVotesByParticipant(dbpool *pgxpool.Pool, ctx context.Context) (map[domain.Participant]int, error) {
	query := `select
			P.participant_id,
			P.participant_name,
			T.votes
		from
			participants as P
		inner join (
			select
				participant_id,
				count(*) as votes
			from
				votes
			group by
				participant_id) as T
				on
			T.participant_id = P.participant_id;`
	return getManyFromQuery(dbpool, ctx, query)
}

func getManyFromQuery(dbpool *pgxpool.Pool, ctx context.Context, query string) (map[domain.Participant]int, error) {
	rows, err := dbpool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query participants: %w", err)
	}
	defer rows.Close()

	result := map[domain.Participant]int{}

	for rows.Next() {
		var participant domain.Participant
		var votes int
		err := rows.Scan(&participant.ParticipantID, &participant.Name, &votes)
		participant.Name = strings.TrimSpace(participant.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan participant: %w", err)
		}
		result[participant] = votes
	}

	return result, nil
}
