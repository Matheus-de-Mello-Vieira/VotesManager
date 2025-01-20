package redisdatamapper

import (
	"bbb-voting/voting-commons/domain"
	"cmp"
	"context"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type VoteDataMapperRedisDecorator struct {
	redis redis.Client
	base  domain.VoteRepository
	participantRepository domain.ParticipantRepository
}

func DecorateVoteDataRepository(base domain.VoteRepository, redis redis.Client, participantRepository domain.ParticipantRepository) VoteDataMapperRedisDecorator {
	return VoteDataMapperRedisDecorator{redis, base, participantRepository}
}

func (mapper VoteDataMapperRedisDecorator) SaveOne(ctx context.Context, vote *domain.Vote) error {
	pipeline := mapper.redis.TxPipeline()

	mapper.saveOneRedis(ctx, vote, pipeline)

	_, err := pipeline.Exec(ctx)
	if err != nil {
		return err
	}

	return mapper.base.SaveOne(ctx, vote)
}

func (mapper VoteDataMapperRedisDecorator) SaveMany(ctx context.Context, votes []domain.Vote) error {
	pipeline := mapper.redis.TxPipeline()

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

	hourKey := fmt.Sprint(vote.GetHour().Unix())
	pipeline.HIncrBy(ctx, "votes:by:hour", hourKey, 1)

	return nil
}

func (mapper VoteDataMapperRedisDecorator) GetGeneralTotal(ctx context.Context) (int, error) {
	resultStr, err := mapper.redis.Get(ctx, "votes:total").Result()
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(resultStr)
}
func (mapper VoteDataMapperRedisDecorator) GetTotalByParticipant(ctx context.Context) (map[domain.Participant]int, error) {
	participantsIdByVotes, err := mapper.redis.HGetAll(ctx, "votes:by:participant").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get votes by participant: %w", err)
	}

	result := map[domain.Participant]int{}
	for participantIdStr, voteStr := range participantsIdByVotes {
		participantId, err := strconv.Atoi(participantIdStr)
		if err != nil {
			return nil, err
		}
		vote, err := strconv.Atoi(voteStr)
		if err != nil {
			return nil, err
		}

		participant, err := mapper.participantRepository.FindByID(ctx, participantId)
		if err != nil {
			return nil, err
		}
		result[*participant] = vote
	}

	return result, nil
}

func (mapper VoteDataMapperRedisDecorator) GetTotalByHour(ctx context.Context) ([]domain.TotalByHour, error) {
	results, err := mapper.redis.HGetAll(ctx, "votes:by:hour").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get votes by hour: %w", err)
	}

	result := []domain.TotalByHour{}
	for hourStr, totalStr := range results {
		hour, err := strconv.ParseInt(hourStr, 10, 64)
		if err != nil {
			return nil, err
		}
		total, err := strconv.Atoi(totalStr)
		if err != nil {
			return nil, err
		}

		element := domain.TotalByHour{
			Total: total,
			Hour:  time.Unix(hour, 0),
		}
		result = append(result, element)
	}

	slices.SortFunc(result, func(a, b domain.TotalByHour) int {
		return cmp.Compare(a.Hour.Unix(), b.Hour.Unix())
	})

	return result, nil
}
