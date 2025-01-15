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

func DecorateVoteDataMapperWithRedis(base domain.VoteRepository, redis redis.Client) VoteDataMapperRedisDecorator {
	return VoteDataMapperRedisDecorator{redis, base}
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
	pipeline.Incr(ctx, "votes:total")

	participantCountKey := fmt.Sprint(vote.Participant.ParticipantID)
	pipeline.HIncrBy(ctx, "votes:by:participant", participantCountKey, 1)

	hourKey := vote.Timestamp.Format("15")
	pipeline.HIncrBy(ctx, "votes:by:hour", hourKey, 1)

	return nil
}
