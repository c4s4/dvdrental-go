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
	Port     = "8080"
	Indent = "  "
	SqlActor = `
		SELECT
        	actor_id,
			first_name,
			last_name
		FROM actor
		WHERE actor_id = $1;`
	SqlFilm = `
		SELECT
			film_id,
			title,
			release_year
		FROM film
		WHERE film_id = $1;
	`
	SqlFilmsWithActor = `
		SELECT
			film.film_id,
			film.title,
			film.release_year
		FROM film, film_actor, actor
		WHERE film_actor.film_id = film.film_id AND
			film_actor.actor_id = actor.actor_id AND
			film_actor.actor_id = $1;
	`
)

type Actor struct {
	Id        int
	FirstName string
	LastName  string
}

type Film struct {
	Id    int
	Title string
	Year  int
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
	out, err := json.MarshalIndent(actor, "", Indent)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.String(http.StatusOK, string(out))
}

func film(ctx *gin.Context) {
	id := ctx.Param("id")
	row := db.QueryRow(SqlFilm, id)
	film := Film{}
	err := row.Scan(
		&film.Id,
		&film.Title,
		&film.Year)
	if err != nil {
		ctx.String(http.StatusNotFound, err.Error())
		return
	}
	out, err := json.MarshalIndent(film, "", Indent)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.String(http.StatusOK, string(out))
}

func filmsWithActor(ctx *gin.Context) {
	actorID := ctx.Param("actor_id")
	rows, err := db.Query(SqlFilmsWithActor, actorID)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()
	films := []Film{}
	for rows.Next() {
		film := Film{}
		err := rows.Scan(
			&film.Id,
			&film.Title,
			&film.Year)
		if err != nil {
			ctx.String(http.StatusNotFound, err.Error())
			return
		}
		films = append(films, film)
	}
	out, err := json.MarshalIndent(films, "", Indent)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.String(http.StatusOK, string(out))
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/actor/:id", actor)
	router.GET("/film/:id", film)
	router.GET("/films/actor/:actor_id", filmsWithActor)
	return router
}

func main() {
	err := ConnectDb()
	if err != nil {
		panic(err)
	}
	log.Print("Connected to database")
	defer db.Close()
	gin.SetMode(gin.ReleaseMode)
	router := setupRouter()
	log.Print("Listening port " + Port)
	router.Run()
	log.Fatal(http.ListenAndServe("localhost:"+Port, nil))
}
