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
	err := s.repo.Set(ctx, item.Email, item.Code, 10*time.Minute).Err()
	//email => otp
	// key val pair

	if err != nil {
		return err
	}
	return nil
}

func (g *Repository)GetOTP(email string) (string, error){
	ctx := context.Background()
	val, err := g.repo.Get(ctx, email).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (d *Repository)DeleteOTP(email string) error{
	ctx := context.Background()
	err := d.repo.Del(ctx, email).Err()
	if err != nil {
		return err
	}
	return nil
}