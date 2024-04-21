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
	err := rdb.Set(context.Background(), "foo", "bar", 0).Err()
	if err != nil {
		fmt.Println("SET Fail", err)
		t.Fail()
	}

	fmt.Println("SET PASS")
	val, err := rdb.Get(context.Background(), "foo").Result()
	if err != nil && val != "bar" {
		fmt.Println("GET Fail", err)
		t.Fail()
	}
	fmt.Println("GET PASS")
}
