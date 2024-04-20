package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestNewRedisClinet(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:3000",
		Password: "",
		DB:       0,
	})
	fmt.Println("TestNewRedisClinet run")
	fmt.Println(rdb)
	err := rdb.Set(context.Background(), "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}
}
