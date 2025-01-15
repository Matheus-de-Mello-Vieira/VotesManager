package redisdatamapper

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type ParticipantDataMapperRedisDecorator struct {
	redis            redis.Client
	base             domain.ParticipantRepository
	participantsById map[int]domain.Participant
	participants     []domain.Participant
}

func DecorateParticipantDataMapperWithRedis(base domain.ParticipantRepository, redis redis.Client) ParticipantDataMapperRedisDecorator {
	return ParticipantDataMapperRedisDecorator{redis, base, nil, nil}
}

func (mapper ParticipantDataMapperRedisDecorator) FindAll(ctx context.Context) ([]domain.Participant, error) {
	var err error
	if mapper.participants == nil {
		mapper.participants, err = mapper.base.FindAll(ctx)
		if err != nil {
			return nil, err
		}
	}
	return mapper.participants, nil
}
func (mapper ParticipantDataMapperRedisDecorator) FindByID(ctx context.Context, id int) (*domain.Participant, error) {
	participantById, err := mapper.getParticipantById(ctx)
	if err != nil {
		return nil, err
	}
	participant := participantById[id]
	return &participant, nil
}
func (mapper ParticipantDataMapperRedisDecorator) getParticipantById(ctx context.Context) (map[int]domain.Participant, error) {
	var err error
	if mapper.participantsById == nil {
		mapper.participantsById, err = mapper.fetchParticipantById(ctx)
		if err != nil {
			return nil, err
		}
	}
	return mapper.participantsById, nil
}
func (mapper ParticipantDataMapperRedisDecorator) fetchParticipantById(ctx context.Context) (map[int]domain.Participant, error) {
	participants, err := mapper.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := map[int]domain.Participant{}
	for _, participant := range participants {
		result[participant.ParticipantID] = participant
	}

	return result, nil
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
	resultStr, err := mapper.redis.Get(ctx, "votes:total").Result()
	if err != nil {
		return nil, err
	}
	result, err := strconv.Atoi(resultStr)
	return &result, err
}
func (mapper ParticipantDataMapperRedisDecorator) getVotesByParticipant(ctx context.Context) (map[domain.Participant]int, error) {
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

		participant, err := mapper.FindByID(ctx, participantId)
		if err != nil {
			return nil, err
		}
		result[*participant] = vote
	}

	return result, nil
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
