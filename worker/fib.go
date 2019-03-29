package main

import (
	"github.com/go-redis/redis"
	"log"
	"strconv"
	"sync"
	"time"
)

func fib(idx int) int {
	if idx < 2 {
		return 1
	}
	return fib(idx-1) + fib(idx-2)
}

func updateRedisOnMsg(client *redis.Client, msgCh <-chan *redis.Message) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		log.Println("waiting for messages...")
		for msg := range msgCh {
			number, err := strconv.Atoi(msg.Payload)
			if err != nil {
				log.Printf("error converting payload %s to number: %s", msg.Payload, err)
				continue
			}
			fibVal := fib(number)
			start := time.Now()
			client.HSet("values", msg.Payload, fibVal)
			elapsed := time.Since(start)
			log.Printf("fibonacci of %d set in %v", number, elapsed)
		}
		log.Fatal("messages channel closed!")
	}()
	wg.Wait()
}

func processRedisMsg(client *redis.Client) {
	pubSub := client.Subscribe("insert")
	subCh := pubSub.Channel()
	updateRedisOnMsg(client, subCh)
}
