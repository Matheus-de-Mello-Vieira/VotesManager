package redisdatamapper

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const lockValue = "locked"

type RedisLock struct {
	client *redis.Client
	key    string
	ttl    time.Duration
}

func NewRedisLock(client *redis.Client, key string, ttl time.Duration) *RedisLock {
	return &RedisLock{
		client: client,
		key:    key,
		ttl:    ttl,
	}
}

func (lock *RedisLock) Lock(ctx context.Context) (bool, error) {
	// Use SET with NX and EX options to set the lock atomically
	success, err := lock.client.SetNX(ctx, lock.key, lockValue, lock.ttl).Result()
	if err != nil {
		return false, err
	}
	return success, nil
}

func (lock *RedisLock) Unlock(ctx context.Context) error {
	script := `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end`
	_, err := lock.client.Eval(ctx, script, []string{lock.key}, lockValue).Result()
	return err
}
