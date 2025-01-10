package postgresqldatamapper

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VoteDataMapper struct {
	connector PostgresqlConnector
}

func NewVoteDataMapper(connector PostgresqlConnector) VoteDataMapper {
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

func getGeneralTotal(dbpool *pgxpool.Pool, ctx context.Context) (*int, error) {
	query := "select count(*) as votes from votes"

	var result int
	err := dbpool.QueryRow(ctx, query).Scan(&result)
	if err != nil {
		result = -1
		return &result, fmt.Errorf("failed to get genetal total: %w", err)
	}

	return &result, nil
}

func getVotesByHour(dbpool *pgxpool.Pool, ctx context.Context) ([]domain.TotalByHour, error) {
	query := `select
			date_part('hour', timestamp) :: integer,
			count(*) as votes 
		from
			votes
		group by
			DATE_PART('hour', timestamp) :: integer
		order by 
			DATE_PART('hour', timestamp) :: integer`

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
