package main

import (
	"fmt"
	"github.com/briferz/complexapp/shared/keys"
	"github.com/go-redis/redis"
	"log"
	"os"
	"os/signal"
)

func main() {
	opt := redis.Options{
		Addr:     keys.RedisAddr(),
		Password: keys.RedisPass(),
	}
	client := redis.NewClient(&opt)

	err := client.Ping().Err()
	if err != nil {
		log.Fatal("unable to reach Redis: ", err)
	}

	log.Println("success reaching Redis.")

	processRedisMsg(client)
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh,os.Interrupt,os.Kill)
	sig := <-signalCh
	fmt.Printf("Received signal %s. Exiting...\n", sig)
}

