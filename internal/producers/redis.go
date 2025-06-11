package producers

import (
	"context"
	"fmt"
	"hermesx/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient interface {
	Close() error
	Ping(ctx context.Context) *redis.StatusCmd
	Get(ctx context.Context, key string) (string, error)
	Incr(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, ttl time.Duration) *redis.BoolCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error)
	Del(ctx context.Context, key string) *redis.IntCmd
}

type Redis struct {
	store *redis.Client
}

func NewRedis(cfg *config.RedisConfig) RedisClient {
	addr := fmt.Sprintf("%s:%v", cfg.Host, cfg.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password, // no password set
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize, // use default DB
	})
	// Initialize custom config
	// store := redis.New(redis.Config{
	// 	Host:      cfg.Host,
	// 	Port:      cfg.Port,
	// 	Password:  cfg.Password,
	// 	Database:  cfg.DB,
	// 	Reset:     false,
	// 	TLSConfig: nil,
	// 	PoolSize:  cfg.PoolSize,
	// })
	return &Redis{
		store: rdb,
	}
}
func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.store.Get(ctx, key).Result()
}
func (r *Redis) Incr(ctx context.Context, key string) (int64, error) {
	return r.store.Incr(ctx, key).Result()
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	return r.store.Set(ctx, key, value, expiration).Result()
}

func (r *Redis) Close() error {
	return r.store.Close()
}

func (r *Redis) Ping(ctx context.Context) *redis.StatusCmd {
	return r.store.Ping(ctx)
}

func (r *Redis) Expire(ctx context.Context, key string, ttl time.Duration) *redis.BoolCmd {
	return r.store.Expire(ctx, key, ttl)
}

func (r *Redis) Del(ctx context.Context, key string) *redis.IntCmd {
	return r.store.Del(ctx, key)
}
