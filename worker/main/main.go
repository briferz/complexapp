package main

import (
	"github.com/go-redis/redis"
	"log"
)

func main() {
	opt := redis.Options{
		Addr:               "localhost:6379",
		Password:           "",
	}
	client := redis.NewClient(&opt)

	err := client.Ping().Err()
	if err != nil {
		log.Fatal("unable to reach Redis: ",err)
	}

	log.Println("success reaching Redis.")

}
