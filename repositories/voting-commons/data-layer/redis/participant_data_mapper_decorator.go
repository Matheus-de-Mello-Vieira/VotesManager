package redisdatamapper

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type ParticipantDataMapperRedisDecorator struct {
	redis    *redis.Client
	base     domain.ParticipantRepository
	baseLock *RedisLock
	ttl      time.Duration
}

const participantsKey = "participants"

func DecorateParticipantRepository(base domain.ParticipantRepository, redis *redis.Client, ttl time.Duration) *ParticipantDataMapperRedisDecorator {
	lock_ttl, _ := time.ParseDuration("30s")

	lock := NewRedisLock(redis, "participants:lock", lock_ttl)
	return &ParticipantDataMapperRedisDecorator{redis: redis, base: base, baseLock: lock, ttl: ttl}
}

func (mapper ParticipantDataMapperRedisDecorator) FindAll(ctx context.Context) ([]domain.Participant, error) {
	err := mapper.loadCacheIfHaveNotLoaded(ctx)
	if err != nil {
		return nil, err
	}

	resultStr, err := mapper.redis.Get(ctx, participantsKey).Result()
	if err != nil {
		return nil, err
	}

	var result []domain.Participant

	json.Unmarshal([]byte(resultStr), &result)

	return result, nil
}
func (mapper ParticipantDataMapperRedisDecorator) FindByID(ctx context.Context, id int) (*domain.Participant, error) {
	participants, err := mapper.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, participant := range participants {
		if participant.ParticipantID == id {
			return &participant, nil
		}
	}

	return nil, nil
}

func (mapper *ParticipantDataMapperRedisDecorator) loadCacheIfHaveNotLoaded(ctx context.Context) error {
	mapper.baseLock.Lock(ctx)
	defer mapper.baseLock.Unlock(ctx)

	isCacheValid, err := mapper.isCacheValid(ctx)
	if err != nil {
		return err
	}

	if !isCacheValid {
		return mapper.loadCache(ctx)
	}

	return nil
}
func (mapper *ParticipantDataMapperRedisDecorator) isCacheValid(ctx context.Context) (bool, error) {
	exists, err := mapper.redis.Exists(ctx, participantsKey).Result()

	if err != nil {
		return false, err
	}

	return exists == 1, nil
}
func (mapper *ParticipantDataMapperRedisDecorator) loadCache(ctx context.Context) error {
	value, err := mapper.fetchCacheTobeValue(ctx)
	if err != nil {
		return err
	}

	err = mapper.redis.Set(ctx, participantsKey, value, mapper.ttl).Err()
	if err != nil {
		return fmt.Errorf("could not set value: %v", err)
	}
	return nil
}

func (mapper *ParticipantDataMapperRedisDecorator) fetchCacheTobeValue(ctx context.Context) (string, error) {
	participants, err := mapper.base.FindAll(ctx)
	if err != nil {
		return "", err
	}

	voteJSON, err := json.Marshal(participants)
	if err != nil {
		return "", fmt.Errorf("failed to serialize vote to JSON: %v", err)
	}

	return string(voteJSON), nil
}
