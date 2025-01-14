package redisdatamapper

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type VoteDataMapperRedisDecorator struct {
	redis redis.Client
	base  domain.VoteRepository
}

func (mapper VoteDataMapperRedisDecorator) SaveOne(ctx context.Context, vote *domain.Vote) error {
	pipeline := mapper.redis.Pipeline()

	mapper.saveOneRedis(ctx, vote, pipeline)

	_, err := pipeline.Exec(ctx)
	if err != nil {
		return err
	}

	return mapper.base.SaveOne(ctx, vote)
}

func (mapper VoteDataMapperRedisDecorator) SaveMany(ctx context.Context, votes []domain.Vote) error {
	pipeline := mapper.redis.Pipeline()

	for _, vote := range votes {
		mapper.saveOneRedis(ctx, &vote, pipeline)
	}

	_, err := pipeline.Exec(ctx)
	if err != nil {
		return err
	}

	return mapper.base.SaveMany(ctx, votes)
}

func (mapper VoteDataMapperRedisDecorator) saveOneRedis(ctx context.Context, vote *domain.Vote, pipeline redis.Pipeliner) error {
	voteKey := fmt.Sprint("vote:", vote.VoteID)
	pipeline.HMSet(ctx, voteKey, map[string]interface{}{
		"id":             vote.VoteID,
		"participant_id": vote.Participant.ParticipantID,
		"timestamp":      vote.Timestamp.Unix(),
	})

	participantCountKey := fmt.Sprint("count:participant:", vote.Participant.ParticipantID)
	pipeline.Incr(ctx, participantCountKey)

	pipeline.Incr(ctx, "count:total")

	hourKey := vote.Timestamp.Format("2006010215")
	pipeline.HIncrBy(ctx, "votes:by:hour", hourKey, 1)

	return nil
}
