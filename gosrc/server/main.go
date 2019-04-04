package main

import (
	"github.com/briferz/complexapp/server/controller"
	"github.com/briferz/complexapp/shared/redisshared"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

func main() {
	log.Printf("Starting...")
	db, err := pgDial()
	if err != nil {
		log.Printf("error: %s", err)
		os.Exit(1)
	}
	defer db.Close()
	log.Println("success reaching database.")

	client, err := redisshared.Client()
	if err != nil {
		log.Print(err)
		os.Exit(2)
	}
	log.Println("success reaching Redis.")

	ctr := controller.New(db, client)

	mux := http.NewServeMux()

	mux.HandleFunc("/values/current", controller.DurationMid(ctr.HandleGetPostgresCurrent))
	mux.HandleFunc("/values/all", controller.DurationMid(ctr.HandleGetRedisValues))
	mux.HandleFunc("/values", controller.DurationMid(ctr.HandlePostValue))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hi!"))
	})

	sigCh := func() chan os.Signal {
		ch := make(chan os.Signal)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			wg.Done()
			signal.Notify(ch, os.Interrupt, os.Kill)
		}()
		wg.Wait()
		return ch
	}()

	go func() {
		log.Println("listening...")
		err = http.ListenAndServe(":8080", mux)
		if err != nil {
			log.Printf("unable to bind server: %s", err)
			os.Exit(3)
		}
	}()

	log.Printf("signaled finish. Signal=%v", <-sigCh)

}
