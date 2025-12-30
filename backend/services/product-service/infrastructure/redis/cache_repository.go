package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheRepository はRedisキャッシュの操作を提供
type CacheRepository struct {
	client *redis.Client
}

// NewCacheRepository はCacheRepositoryの新しいインスタンスを作成
func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{
		client: client,
	}
}

// Get はキャッシュからデータを取得
func (r *CacheRepository) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// Set はキャッシュにデータを保存
func (r *CacheRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, jsonData, expiration).Err()
}

// Delete はキャッシュからデータを削除
func (r *CacheRepository) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// DeletePattern は指定されたパターンに一致するキーを削除
func (r *CacheRepository) DeletePattern(ctx context.Context, pattern string) error {
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := r.client.Del(ctx, iter.Val()).Err(); err != nil {
			return fmt.Errorf("failed to delete key %s: %w", iter.Val(), err)
		}
	}
	return iter.Err()
}

// Ping はRedis接続を確認
func (r *CacheRepository) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
