package main

import (
	"fmt"
	"github.com/briferz/complexapp/shared/redisshared"
	"log"
	"os"
	"os/signal"
)

func main() {
	client, err := redisshared.Client()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("success reaching Redis.")

	processRedisMsg(client)
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt, os.Kill)
	sig := <-signalCh
	fmt.Printf("Received signal %s. Exiting...\n", sig)
}
