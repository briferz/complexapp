package main

import (
	"github.com/briferz/complexapp/server/controller"
	"github.com/briferz/complexapp/shared/redisshared"
	"log"
	"net/http"
)

func main() {

	db, err := pgDial()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	defer db.Close()
	log.Println("success reaching database.")

	client, err := redisshared.Client()
	if err != nil {
		log.Fatal(err)
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

	log.Println("listening...")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("unable to bind server: %s", err)
	}

}
