package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func New(addr, password string) *Redis {
	client := redis.NewClient(&redis.Options{Addr: addr, Password: password})
	return &Redis{Client: client}
}

func (r *Redis) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

func (r *Redis) Close() error {
	return r.Client.Close()
}

func (r *Redis) SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.Client.Set(ctx, key, data, ttl).Err()
}

func (r *Redis) GetJSON(ctx context.Context, key string, dest any) (bool, error) {
	val, err := r.Client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal(val, dest); err != nil {
		return false, err
	}
	return true, nil
}

func (r *Redis) AddBlacklist(ctx context.Context, token string, ttl time.Duration) error {
	return r.Client.Set(ctx, fmt.Sprintf("blacklist:%s", token), "1", ttl).Err()
}

func (r *Redis) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	count, err := r.Client.Exists(ctx, fmt.Sprintf("blacklist:%s", token)).Result()
	return count > 0, err
}

func (r *Redis) SlidingWindowHit(ctx context.Context, key string, limit int64, window time.Duration) (bool, error) {
	now := time.Now().UnixNano()
	pipe := r.Client.TxPipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", time.Now().Add(-window).UnixNano()))
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: fmt.Sprintf("%d", now)})
	pipe.ZCard(ctx, key)
	pipe.Expire(ctx, key, window)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}
	count, err := r.Client.ZCard(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count <= limit, nil
}

