package ratelimiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type rateLimiterRepo interface {
	InsertRequest(email string, totalRequestsSent int) error
	GetRequest(email string) (string, error)

}

type Repository struct{
	repo *redis.Client
}

func NewRepo(rdb *redis.Client) *Repository{
	return &Repository{
		repo: rdb,
	}
}


func (i *Repository) InsertRequest(email string, totalRequestsSent int) error{
	cxt := context.Background()
	key := RedisKeyPrefix + email
	duration, _ := i.repo.TTL(cxt, key).Result()
	if duration == -2 *time.Nanosecond || duration == -1 *time.Nanosecond{
		err := i.repo.Set(cxt, key, totalRequestsSent, TTL).Err()
		return err
	}

	err := i.repo.Set(cxt, key, totalRequestsSent, duration).Err()
	return err
}


func (i *Repository) GetRequest(email string) (string, error){
	cxt := context.Background()
	key := RedisKeyPrefix + email
	val, err := i.repo.Get(cxt, key).Result()
	if err == redis.Nil{
		return "", err
	}
	if err != redis.Nil && err != nil {
		return "", err
	}

	return val, nil
}

