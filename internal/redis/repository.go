package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Repository struct {
	redis *redis.Client
}

type CacheRepository interface{
	Insert(key string, value string) error
	Get(key string) (string, error)
}



func NewRepository (client *redis.Client) *Repository{
	return &Repository{
		redis: client,
	}
}

func (r *Repository) Insert(key string, value string) error {
	ctx := context.Background()
	duration, _ := r.redis.TTL(ctx, key).Result()
	// return of TTL: 
	// type DurationCmd struct {
	// 	baseCmd

	// 	val       time.Duration
	// 	precision time.Duration
	// }

	// -1	The key exists, but it has no associated expire/TTL (it is persistent).
	// -2	The key does not exist.
	if duration == -1 * time.Nanosecond || duration == -2 *time.Nanosecond{
		err := r.redis.Set(ctx, key, value, TTL).Err()
		if err != nil {
			return err
		}
	}

	err := r.redis.Set(ctx, key, value, duration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Get(key string) (string, error) {
	cxt := context.Background()

	val, err := r.redis.Get(cxt, key).Result()
	if err == redis.Nil{
		return "", KeydoesNotExist
	}
	if err != redis.Nil && err != nil {
		return "", err
	}

	return val, nil
}





