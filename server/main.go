package main

import (
	"database/sql"
	"fmt"
	"github.com/briferz/complexapp/shared/keys"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/signal"
)

func main() {

	data, err := keys.PgDataSource()
	if err != nil {
		log.Fatalf("error getting connection data: %s", err)
	}
	db, err := sql.Open("postgres", data)
	if err != nil {
		log.Fatalf("error reaching database: %s", err)
	}
	defer db.Close()

	log.Println("success reaching database.")

	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt, os.Kill)
	sig := <-signalCh
	fmt.Printf("Received signal %s. Exiting...\n", sig)
}
