package main

import (
	"back/platform/snake"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
	"github.com/qkgo/yin"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("CONECTION_STRING"))
	if err != nil {
		log.Fatalln("Connecting to db", err)
	}

	_, err = db.Exec(`
	   	CREATE TABLE IF NOT EXISTS users (
	   		id SERIAL NOT NULL PRIMARY KEY,
	   		username text
	   	);

	   	CREATE TABLE IF NOT EXISTS scores (
	   		idscore SERIAL NOT NULL PRIMARY KEY,
	   		userid INT NOT NULL REFERENCES users(id),
	   		score INT NOT NULL
	   	);
	   	`)
	if err != nil {
		log.Fatalln("create tables error: %w", err)
	}

	object := snake.DatabaseDeclaration(db)

	router := chi.NewRouter()
	router.Use(yin.SimpleLogger)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
	}))

	router.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		res, _ := yin.Event(rw, r)
		items, err := object.Get()
		if err != nil {
			res.SendStatus(http.StatusInternalServerError)
			return
		}
		res.SendJSON(items)
	})

	router.Post("/", func(rw http.ResponseWriter, r *http.Request) {
		res, req := yin.Event(rw, r)
		body := map[string]string{}
		req.BindBody(&body)
		item := snake.Item{
			Username: body["Username"],
			Score:    body["Score"],
		}
		err := object.Add(item)
		if errors.Is(err, snake.ErrInvalidQueryStatement) {
			res.SendStatus(http.StatusInternalServerError)
			return
		}
		if err != nil {
			res.SendStatus(http.StatusInternalServerError)
			return
		}
		err = object.AddScore(item, item.Username)
		if err != nil {
			res.SendStatus(http.StatusInternalServerError)
			return
		}
		res.SetHeader("Access-Control-Allow-Origin", "*")
		res.SendStatus(http.StatusOK)
	})

	router.Put("/", func(rw http.ResponseWriter, r *http.Request) {
		res, req := yin.Event(rw, r)
		body := map[string]string{}
		err := req.BindBody(&body)
		if err != nil {
			res.SendStatus(http.StatusInternalServerError)
			return
		}
		item := snake.Item{
			Id:       body["Id"],
			Username: body["Username"],
			Score:    body["Score"],
		}
		err = object.Put(item.Id, item.Username, item.Score)
		if err != nil {
			res.SendStatus(http.StatusInternalServerError)
			return
		}
		res.SendStatus(http.StatusOK)
	})

	router.Delete("/", func(rw http.ResponseWriter, r *http.Request) {
		res, req := yin.Event(rw, r)
		body := map[string]string{}
		err := req.BindBody(&body)
		if err != nil {
			res.SendStatus(http.StatusInternalServerError)
			return
		}
		item := snake.Item{
			Username: body["Username"],
		}
		err = object.Delete(item.Username)
		if err != nil {
			res.SendStatus(http.StatusInternalServerError)
		}
		res.SendStatus(http.StatusOK)
	})

	err = http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatalln("listening to port 3000: %w", err)
	}
}
