package otp


import (
	"github.com/redis/go-redis/v9"
	"context"
	"time"
)

type Repository struct {
	repo *redis.Client
}

type OTPRepository interface {
	SaveOTP(item OTPItem) error
	GetOTP(email string) (string, error)
 	DeleteOTP(email string) error
}

func NewOTPRepo(client *redis.Client) *Repository{
	return &Repository{
		repo: client,
	}
}

func (s *Repository)SaveOTP(item OTPItem) error{
	ctx := context.Background()
	key := RedisKeyPrefix + item.Email
	err := s.repo.Set(ctx, key, item.Code, 10*time.Minute).Err()
	//email => otp
	// key val pair

	if err != nil {
		return err
	}
	return nil
}

func (g *Repository)GetOTP(email string) (string, error){
	ctx := context.Background()
	key := RedisKeyPrefix + email
	val, err := g.repo.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (d *Repository)DeleteOTP(email string) error{
	ctx := context.Background()
	key := RedisKeyPrefix + email
	err := d.repo.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}