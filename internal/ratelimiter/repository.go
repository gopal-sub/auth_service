package ratelimiter

import (
	"gopal-sub/auth_service/internal/redis"
)

type rateLimiterRepo interface {
	InsertRequest(email string, totalRequestsSent string) error
	GetRequest(email string) (string, error)

}

type Repository struct{
	repo *redis.Repository
}

func NewRepo(repo *redis.Repository) *Repository{
	return &Repository{
		repo: repo,
	}
}


func (r *Repository) InsertRequest(email string, totalRequestsSent string) error{
	key := RedisKeyPrefix + email

	err := r.repo.Insert(key, totalRequestsSent)
	if err != nil {
		return err
	}
	return nil
}


func (r *Repository) GetRequest(email string) (string, error){
	key := RedisKeyPrefix + email

	val, err := r.repo.Get(key)
	if err != nil {
		return "", err
	}
	return val, nil
	
}

