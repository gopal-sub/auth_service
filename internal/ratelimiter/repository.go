package ratelimiter

import (
	"context"
	"strconv"
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

	err = i.repo.Set(cxt, email, totalRequestsSent, time.Duration(totalRequestsSent)).Err()
	

	
	
		
		


	return otpAttemptsExeErr
}


func (i *Repository) GetRequest(email string) (string, error){
	cxt := context.Background()
	val, err := i.repo.Get(cxt, email).Result()
	if err == redis.Nil{
		return "", err
	}
	if err != redis.Nil && err != nil {
		return "", err
	}

	return val, nil
}


// email -> {numberOfReq = }
