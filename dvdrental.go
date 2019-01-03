package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	SqlActor = `
		SELECT
			actor_id,
			first_name,
			last_name
		FROM actor
		WHERE actor_id = $1;`
)

type Actor struct {
	Id        int
	FirstName string
	LastName  string
}

var db *sql.DB

func ConnectDb() error {
	var err error
	config := fmt.Sprintf("host='%s' port='%s' user='%s' password='%s' dbname='%s' sslmode='disable'",
		os.Getenv("DBHOST"),
		os.Getenv("DBPORT"),
		os.Getenv("DBUSER"),
		os.Getenv("DBPASS"),
		os.Getenv("DBNAME"))
	db, err = sql.Open("postgres", config)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func actor(ctx *gin.Context) {
	id := ctx.Param("id")
	row := db.QueryRow(SqlActor, id)
	actor := Actor{}
	err := row.Scan(
		&actor.Id,
		&actor.FirstName,
		&actor.LastName)
	if err != nil {
		ctx.String(http.StatusNotFound, err.Error())
		return
	}
	out, err := json.Marshal(actor)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.String(http.StatusOK, string(out))
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/actor/:id", actor)
	return router
}

func main() {
	err := ConnectDb()
	if err != nil {
		panic(err)
	}
	log.Print("Connected to database")
	defer db.Close()
	router := setupRouter()
	router.Run()
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
