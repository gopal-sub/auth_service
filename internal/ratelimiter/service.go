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


func (r *Service) RateLimit(email string){
	val, err := r.repo.GetRequest(email)

	if err == redis.Nil {
		//insert to redis
		return
	}
	if err != nil {
		//something wrogn with redis
	}
	totalRequestsSent, _ := strconv.Atoi(val)
	if totalRequestsSent < 10 {
		totalRequestsSent += 1
		err = r.repo.InsertRequest(email, totalRequestsSent)
		if err !=nil {
			return
		}
		
		return 
	}


	
	
}



