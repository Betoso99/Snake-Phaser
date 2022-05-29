package main

import (
	"back/platform/snake"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
	"github.com/qkgo/yin"
)

func main() {
	// hostname, error := os.Hostname()
	// print(hostname)
	// if error != nil {
	// 	println("Error getting the host: ", error)
	// }

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
		log.Fatalln("Tables", err)
	}

	object := snake.NewTbl(db)

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
		items := object.Get()
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
		object.Add(item)
		object.AddScore(item, item.Username)
		res.SetHeader("Access-Control-Allow-Origin", "*")
		fmt.Println()
		res.SendStatus(200)

	})

	router.Put("/", func(rw http.ResponseWriter, r *http.Request) {
		res, req := yin.Event(rw, r)
		body := map[string]string{}
		req.BindBody(&body)
		item := snake.Item{
			Id:       body["Id"],
			Username: body["Username"],
			Score:    body["Score"],
		}
		object.Put(item.Id, item.Username, item.Score)
		res.SendStatus(200)
	})

	router.Delete("/", func(rw http.ResponseWriter, r *http.Request) {
		res, req := yin.Event(rw, r)
		body := map[string]string{}
		req.BindBody(&body)
		item := snake.Item{
			Username: body["Username"],
		}
		object.Delete(item.Username)
		res.SendStatus(200)
	})

	http.ListenAndServe(":3000", router)
}
