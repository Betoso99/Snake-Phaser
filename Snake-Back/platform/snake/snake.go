package snake

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

type Record struct {
	DB *sql.DB
}

func NewTbl(db *sql.DB) *Record {

	/* stmt, _ := db.Prepare(`
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
	stmt.Exec() */
	return &Record{
		DB: db,
	}
}

func (record *Record) getid(username string) int {
	var value int
	stmt, err := record.DB.Prepare(`
		SELECT id FROM users WHERE username = $1;
	`)
	if err != nil {
		fmt.Println("Error:", err)
	}
	row, err2 := stmt.Query(username)
	if err != nil {
		fmt.Println("Error:", err2)
	}
	for row.Next() {
		row.Scan(&value)
	}
	return value
}

func (record *Record) Add(item Item) {
	stmt, err := record.DB.Prepare(`
		INSERT INTO users(username) values($1);
	`)
	if err != nil {
		fmt.Println("Error:", err)
	}
	stmt.Exec(item.Username)
}

func (record *Record) AddScore(item Item, username string) {
	stmt, err := record.DB.Prepare(`
		INSERT INTO scores(userid, score) values ($1, $2);
	`)
	if err != nil {
		fmt.Println("Error:", err)
	}
	tempID := record.getid(username)
	num, err := strconv.Atoi(item.Score)
	stmt.Exec(tempID, num)
}

func (record *Record) Delete(username string) {
	stmt, err := record.DB.Prepare(`
		DELETE FROM users WHERE id = $1;
	`)
	if err != nil {
		fmt.Println("Error:", err)
	}
	tempID := record.getid(username)
	record.DeleteScore(tempID)
	stmt.Exec(tempID)
}

func (record *Record) DeleteScore(id int) {
	stmt, err := record.DB.Prepare(`
		DELETE FROM scores WHERE userid = $1;
	`)
	if err != nil {
		fmt.Println("Error:", err)
	}
	stmt.Exec(id)
}

func (record *Record) Put(id string, user string, score string) {
	stmt1, err1 := record.DB.Prepare(`
		UPDATE users SET username = $1 WHERE id = $2;
	`)
	if err1 != nil {
		fmt.Println("Error:", err1)
	}
	stmt2, err2 := record.DB.Prepare(`
		UPDATE scores SET score = $1 WHERE userid = $2;
	`)
	if err2 != nil {
		fmt.Println("Error:", err2)
	}

	stmt1.Exec(user, id)
	stmt2.Exec(score, id)
}

func (record *Record) Get() []Item {
	var username string
	var score string
	var id string
	items := []Item{}
	rows, err := record.DB.Query(`
		SELECT u.id, u.username, s.score FROM users AS u INNER JOIN scores AS s ON u.id = s.userid ORDER BY s.score desc LIMIT 10;
	`)
	if err != nil {
		log.Fatalln("------------------------------------------------------------------------", err)
	}

	for rows.Next() {
		rows.Scan(&id, &username, &score)
		item := Item{
			Id:       id,
			Username: username,
			Score:    score,
		}
		items = append(items, item)
	}
	return items
}
