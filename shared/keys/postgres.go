package keys

import (
	"errors"
	"fmt"
	"os"
)

const (
	envPgUser     = "PGUSER"
	envPgHost     = "PGHOST"
	envPgDatabase = "PGDATABASE"
	envPgPassword = "PGPASSWORD"
	envPgPort     = "PGPORT"
)

func PgUser() (string, bool) {
	return os.LookupEnv(envPgUser)
}

func PgHost() (string, bool) {
	return os.LookupEnv(envPgHost)
}

func PgDatabase() (string, bool) {
	return os.LookupEnv(envPgDatabase)
}

func PgPassword() (string, bool) {
	return os.LookupEnv(envPgPassword)
}

func PgPort() (string, bool) {
	return os.LookupEnv(envPgPort)
}

func PgDataSource() (string, error) {

	host, ok := PgHost()
	if !ok {
		return "", errors.New("no Postgres Host configured")
	}

	port, ok := PgPort()
	if !ok {
		return "", errors.New("no Postgres Port configured")
	}

	user, ok := PgUser()
	if !ok {
		return "", errors.New("no Postgres User configured")
	}

	pass, ok := PgPassword()
	if !ok {
		return "", errors.New("no Postgres Password configured")
	}

	db, ok := PgDatabase()
	if !ok {
		return "", errors.New("no Postgres Database configured")
	}

	return fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, pass, db), nil
}
