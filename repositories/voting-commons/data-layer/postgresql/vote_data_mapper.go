package postgresqldatamapper

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VoteDataMapper struct {
	connector *PostgresqlConnector
}

func NewVoteDataMapper(connector *PostgresqlConnector) VoteDataMapper {
	return VoteDataMapper{connector}
}

func (mapper VoteDataMapper) SaveOne(ctx context.Context, vote *domain.Vote) error {
	dbpool, err := mapper.connector.openConnection(ctx)
	if err != nil {
		return err
	}
	defer dbpool.Close()

	query := `
		INSERT INTO votes (participant_id, timestamp)
		VALUES ($1, $2) RETURNING vote_id
	`

	err = dbpool.QueryRow(ctx, query, (*vote).Participant.ParticipantID, vote.Timestamp).Scan(vote.VoteID)
	if err != nil {
		return fmt.Errorf("failed to insert vote: %w", err)
	}

	return nil
}

func (mapper VoteDataMapper) SaveMany(ctx context.Context, votes []domain.Vote) error {
	dbpool, err := mapper.connector.openConnection(ctx)
	if err != nil {
		return err
	}
	defer dbpool.Close()

	// Create a new batch
	batch := &pgx.Batch{}

	// Queue an INSERT for each vote
	for _, v := range votes {
		batch.Queue(
			`INSERT INTO votes (participant_id, timestamp) 
             VALUES ($1, $2) 
             RETURNING vote_id`,
			v.Participant.ParticipantID,
			v.Timestamp,
		)
	}

	// Send the batch to the database
	br := dbpool.SendBatch(ctx, batch)
	defer br.Close()

	// Collect the generated vote_ids for each inserted record
	for i := range votes {
		err := br.QueryRow().Scan(&votes[i].VoteID)
		if err != nil {
			return fmt.Errorf("failed to insert vote at index %d: %w", i, err)
		}
	}

	return nil

}

func (mapper VoteDataMapper) GetGeneralTotal(ctx context.Context) (int, error) {
	dbpool, err := mapper.connector.openConnection(ctx)
	if err != nil {
		return -1, err
	}
	defer dbpool.Close()

	query := "select count(*) as votes from votes"

	var result int
	err = dbpool.QueryRow(ctx, query).Scan(&result)
	if err != nil {
		return -1, fmt.Errorf("failed to get genetal total: %w", err)
	}

	return result, nil
}

func (mapper VoteDataMapper) GetTotalByHour(ctx context.Context) ([]domain.TotalByHour, error) {
	dbpool, err := mapper.connector.openConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer dbpool.Close()

	query := `select hourTimestamp, count(*) as votes
		from (select to_timestamp(FLOOR(extract(epoch from timestamp) / (60 * 60)) * 60 * 60) as hourTimestamp from votes) as T
		group by hourTimestamp
		order by hourTimestamp`

	rows, err := dbpool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query votes per hour: %w", err)
	}
	defer rows.Close()

	var result []domain.TotalByHour
	for rows.Next() {
		var totalByHour domain.TotalByHour
		err := rows.Scan(&totalByHour.Hour, &totalByHour.Total)

		if err != nil {
			return nil, fmt.Errorf("failed to scan participant: %w", err)
		}
		result = append(result, totalByHour)
	}

	return result, nil
}
func (mapper VoteDataMapper) GetTotalByParticipant(ctx context.Context) (map[domain.Participant]int, error) {
	dbpool, err := mapper.connector.openConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer dbpool.Close()

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
