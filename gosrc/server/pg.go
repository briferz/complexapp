package main

import (
	"database/sql"
	"fmt"
	"github.com/briferz/complexapp/shared/keys"
	_ "github.com/lib/pq"
	"log"
)

func pgDial() (*sql.DB, error) {
	data, err := keys.PgDataSource()
	if err != nil {
		log.Fatalf("error getting connection data: %s", err)
	}
	db, err := sql.Open("postgres", data)
	if err != nil {
		return nil, fmt.Errorf("error reaching database: %s", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error on ping: %s", err)
	}

	err = pgSetup(db)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func pgSetup(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS values(number int)")
	if err != nil {
		return fmt.Errorf("creating values table: %s", err)
	}
	return nil
}
