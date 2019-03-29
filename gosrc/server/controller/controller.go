package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func DurationMid(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handlerFunc(w, r)
		log.Printf("%v from %s", time.Since(start), r.RemoteAddr)
	}
}

type controller struct {
	db     *sql.DB
	client *redis.Client
}

func New(dbClient *sql.DB, redisClient *redis.Client) *controller {
	return &controller{
		db:     dbClient,
		client: redisClient,
	}
}

func (c *controller) HandlePostValue(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodOptions:
		cors(w, r, http.MethodPost)
	case http.MethodPost:
		c.handlePostValue(w, r)
	default:
		http.NotFound(w, r)
	}

}

func (c *controller) handlePostValue(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading request body: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	parseStr := &struct {
		Index *int `json:"index"`
	}{}

	err = json.Unmarshal(body, parseStr)
	if err != nil {
		log.Printf("error parsing request body: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if parseStr.Index == nil {
		http.Error(w, "Index is null", 422)
		return
	}

	index := *parseStr.Index

	if index > 40 {
		http.Error(w, "Index is too high", 422)
		return
	}

	err = c.setRedisNewValue(index)
	if err != nil {
		log.Print("Redis failed with error: ", err)
		http.Error(w, http.StatusText(http.StatusFailedDependency), http.StatusFailedDependency)
		return
	}

	err = c.setPgNewValue(index)
	if err != nil {
		log.Print("Postgres failed with error: ", err)
		http.Error(w, http.StatusText(http.StatusFailedDependency), http.StatusFailedDependency)
		return
	}

	err = writeJSON(w, map[string]interface{}{"working": true})
	if err != nil {
		log.Print("failed to write response body: ", err)
	}
}

func (c *controller) setPgNewValue(index int) error {
	_, err := c.db.Exec("INSERT INTO values(number) VALUES($1)", index)
	return err
}

func (c *controller) setRedisNewValue(index int) error {
	pipe := c.client.TxPipeline()
	defer pipe.Close()

	pipe.HSet("values", strconv.Itoa(index), "Not yet defined")
	pipe.Publish("insert", index)

	_, err := pipe.Exec()
	return err
}

func (c *controller) HandleGetRedisValues(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
	cmd := c.client.HGetAll("values")
	if err := cmd.Err(); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err := writeJSON(w, cmd.Val())
	if err != nil {
		log.Printf("error writing response: %s", err)
	}
}

type connEmployer func(w http.ResponseWriter, r *http.Request, conn *sql.Conn)

func (c *controller) withConn(w http.ResponseWriter, r *http.Request, employer connEmployer) {
	conn, err := c.db.Conn(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusFailedDependency), http.StatusFailedDependency)
		return
	}
	defer conn.Close()

	employer(w, r, conn)
}

func (c *controller) HandleGetPostgresCurrent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
	c.withConn(w, r, getAll)
}

func getAll(w http.ResponseWriter, r *http.Request, conn *sql.Conn) {
	row, err := conn.QueryContext(r.Context(), "SELECT * FROM values")
	if err != nil {
		http.Error(w, fmt.Sprintf("query error (%s)", err), http.StatusInternalServerError)
	}
	defer row.Close()

	rows := make([]int, 0)
	for row.Next() {
		var newInt int
		err := row.Scan(&newInt)
		if err != nil {
			http.Error(w, fmt.Sprintf("row retrieval error (%s)", err), http.StatusInternalServerError)
			return
		}
		rows = append(rows, newInt)
	}
	if err := row.Err(); err != nil {
		http.Error(w, fmt.Sprintf("row advancing error (%s)", err), http.StatusInternalServerError)
		return
	}
	err = writeJSON(w, rows)
	if err != nil {
		log.Printf("error writing response: %s", err)
	}
}

func cors(w http.ResponseWriter, r *http.Request, methods ...string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ", "))
}

func writeJSON(w http.ResponseWriter, v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	w.Header().Set("Content-type", "application/json")
	_, err = w.Write(bytes)
	return err
}
