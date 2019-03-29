package redisshared

import (
	"fmt"
	"github.com/briferz/complexapp/shared/keys"
	"github.com/go-redis/redis"
)

func Client() (*redis.Client, error) {
	opt := redis.Options{
		Addr:     keys.RedisAddr(),
		Password: keys.RedisPass(),
	}
	client := redis.NewClient(&opt)

	err := client.Ping().Err()
	if err != nil {
		return nil, fmt.Errorf("unable to reach Redis: %s", err)
	}

	return client, nil
}
