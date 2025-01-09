package postgresqldatamapper

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"fmt"
)

type PostgresqlVoteDataMapper struct {
	connector PostgresqlConnector
}

func (mapper *PostgresqlVoteDataMapper) SaveOne(ctx context.Context, vote *domain.Vote) error {
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
