package postgresqldatamapper

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"fmt"
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
