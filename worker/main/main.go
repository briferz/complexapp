package main

import (
	"github.com/go-redis/redis"
	"log"
)

func main() {
	opt := redis.Options{
		Addr:               "localhost",
		Password:           "",
	}
	client := redis.NewClient(&opt)

	err := client.Ping().Err()
	if err != nil {
		log.Fatal("unable to reach Redis.")
	}

	log.Println("success reaching Redis.")

}
