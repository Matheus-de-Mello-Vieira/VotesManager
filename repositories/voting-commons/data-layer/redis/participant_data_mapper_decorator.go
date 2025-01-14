package redisdatamapper

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type ParticipantDataMapperRedisDecorator struct {
	redis redis.Client
	base  domain.ParticipantRepository
}

func DecorateParticipantDataMapperWithRedis(base domain.ParticipantRepository, redis redis.Client) ParticipantDataMapperRedisDecorator {
	return ParticipantDataMapperRedisDecorator{redis, base}
}

func (mapper ParticipantDataMapperRedisDecorator) FindAll(ctx context.Context) ([]domain.Participant, error) {
	return mapper.base.FindAll(ctx)
}
func (mapper ParticipantDataMapperRedisDecorator) FindByID(ctx context.Context, id int) (*domain.Participant, error) {
	return mapper.base.FindByID(ctx, id)
}

func (mapper ParticipantDataMapperRedisDecorator) GetRoughTotals(ctx context.Context) (map[domain.Participant]int, error) {
	return mapper.getVotesByParticipant(ctx)
}

func (mapper ParticipantDataMapperRedisDecorator) GetThoroughTotals(ctx context.Context) (*domain.ThoroughTotals, error) {
	generalTotal, err := mapper.getGeneralTotal(ctx)
	if err != nil {
		return nil, err
	}

	totalByParticipant, err := mapper.getVotesByParticipant(ctx)
	if err != nil {
		return nil, err
	}

	totalByHour, err := mapper.getVotesByHour(ctx)
	if err != nil {
		return nil, err
	}

	result := domain.ThoroughTotals{GeneralTotal: *generalTotal, TotalByHour: totalByHour, TotalByParticipant: totalByParticipant}

	return &result, nil
}

func (mapper ParticipantDataMapperRedisDecorator) getGeneralTotal(ctx context.Context) (*int, error) {
	resultStr, err := mapper.redis.Get(ctx, "count:total").Result()
	if err != nil {
		return nil, err
	}
	result, err := strconv.Atoi(resultStr)
	return &result, err
}
func (mapper ParticipantDataMapperRedisDecorator) getVotesByParticipant(ctx context.Context) (map[domain.Participant]int, error) {
	participants, err := mapper.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := map[domain.Participant]int{}
	for _, participant := range participants {
		count, err := mapper.getCountOfParticipant(ctx, participant.ParticipantID)

		if err != nil {
			return nil, err
		}

		result[participant] = *count
	}

	return result, nil
}
func (mapper ParticipantDataMapperRedisDecorator) getCountOfParticipant(ctx context.Context, participantId int) (*int, error) {
	key := fmt.Sprint("count:participant:", participantId)

	countStr, err := mapper.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil { // don't exists
			count := 0
			return &count, nil
		}
		return nil, err
	}
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return nil, err
	}

	return &count, nil
}

func (mapper ParticipantDataMapperRedisDecorator) getVotesByHour(ctx context.Context) ([]domain.TotalByHour, error) {
	results, err := mapper.redis.HGetAll(ctx, "votes:by:hour").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get votes by hour: %w", err)
	}

	result := []domain.TotalByHour{}
	for hourStr, totalStr := range results {
		hour, err := strconv.Atoi(hourStr)
		if err != nil {
			return nil, err
		}
		total, err := strconv.Atoi(totalStr)
		if err != nil {
			return nil, err
		}

		element := domain.TotalByHour{
			Total: total,
			Hour:  hour,
		}
		result = append(result, element)
	}

	return result, nil
}
