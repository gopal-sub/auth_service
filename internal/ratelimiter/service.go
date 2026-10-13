package ratelimiter

import (
	"strconv"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	repo rateLimiterRepo
}

func NewService(repo rateLimiterRepo) *Service {
	return &Service{
		repo: repo,
	}
}


func (r *Service) RateLimit(email string) error{
	val, err := r.repo.GetRequest(email)
	if err == redis.Nil {
		//insert to redis
		// its a fresh insert
		err = r.repo.InsertRequest(email, "1")
		if err != nil {
			return err
		}
		return nil
	}
	if err != nil {
		//something wrong with redis
		return err
	}
	totalRequestsSent, _ := strconv.Atoi(val)
	if totalRequestsSent < MaxReqPerTTL {
		totalRequestsSent += 1
		err = r.repo.InsertRequest(email, string(totalRequestsSent))
		if err != nil {
			return err
		}
		
		return nil
	}

	return OTPAttemptsExeErr
}



