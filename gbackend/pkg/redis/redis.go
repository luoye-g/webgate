package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

var rdb *goredis.Client

func InitRedis(host, port, pass string) {
	rdb = goredis.NewClient(&goredis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
		DB:   0,
	})

	if err := rdb.Set(context.Background(), "init_key", "0", time.Second).Err(); err != nil {
		panic(err)
	}

	val, err := rdb.Get(context.Background(), "init_key").Result()
	if err != nil {
		panic(err)
	}

	if val != "0" {
		panic("init redis error")
	}
}

func GetRedisCli() *goredis.Client {
	return rdb
}
