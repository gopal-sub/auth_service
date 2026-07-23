package redis


import (
	"github.com/redis/go-redis/v9"
    "context"
)


func New() (*redis.Client, error){
    ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })
    if err := rdb.Ping(ctx).Err(); err !=nil {
        return nil, err
    }
	return rdb, nil
}