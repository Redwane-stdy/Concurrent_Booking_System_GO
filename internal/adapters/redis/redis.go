package redis

import (
	"context"
	"fmt"
	"log"

	goredis "github.com/redis/go-redis/v9"
)

func NewClient(addr string) *goredis.Client {

	rdb := goredis.NewClient(&goredis.Options{
		Addr: addr,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	fmt.Printf("Connected to Redis at %s\n", addr)
	log.Printf("Connected to Redis at %s", addr)
	return rdb
}
