package snake

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrInvalidQueryStatement = errors.New("invalid query statement")
)
var (
	ErrScanCouldntComplete = errors.New("row scan error")
)

type Record struct {
	DB *sql.DB
}

func DatabaseDeclaration(db *sql.DB) *Record {
	return &Record{
		DB: db,
	}
}

func (record *Record) getid(username string) (int, error) {
	var value int
	stmt, err := record.DB.Prepare(`
		SELECT id FROM users WHERE username = $1;
	`)
	if err != nil {
		return 0, ErrInvalidQueryStatement
	}
	row, err := stmt.Query(username)
	if err != nil {
		return 0, fmt.Errorf("find user statement: %w", err)
	}
	for row.Next() {
		row.Scan(&value)
		if err != nil {
			return 0, ErrScanCouldntComplete
		}
	}
	return value, nil
}

func (record *Record) Add(item Item) error {
	stmt, err := record.DB.Prepare(`
		INSERT INTO users(username) values($1);
	`)
	if err != nil {
		return ErrInvalidQueryStatement
	}
	_, err = stmt.Exec(item.Username)
	if err != nil {
		return fmt.Errorf("insert statement error: %w", err)
	}
	return nil
}

func (record *Record) AddScore(item Item, username string) error {
	stmt, err := record.DB.Prepare(`
		INSERT INTO scores(userid, score) values ($1, $2);
	`)
	if err != nil {
		return ErrInvalidQueryStatement
	}
	tempID, err := record.getid(username)
	if err != nil {
		return fmt.Errorf("get id error: %w", err)
	}

	num, err := strconv.Atoi(item.Score)
	if err != nil {
		return fmt.Errorf("return num error: %w", err)
	}
	_, err = stmt.Exec(tempID, num)
	if err != nil {
		return fmt.Errorf("insert statement error: %w", err)
	}
	return nil
}

func (record *Record) Delete(username string) error {
	stmt, err := record.DB.Prepare(`
		DELETE FROM users WHERE id = $1;
	`)
	if err != nil {
		return ErrInvalidQueryStatement
	}
	tempID, err := record.getid(username)
	if err != nil {
		return fmt.Errorf("get id error: %w", err)
	}
	_, err = record.DeleteScore(tempID), nil
	if err != nil {
		return fmt.Errorf("record delete statement error: %w", err)
	}
	_, err = stmt.Exec(tempID)
	if err != nil {
		return fmt.Errorf("delete score statement error: %w", err)
	}
	return nil
}

func (record *Record) DeleteScore(id int) error {
	stmt, err := record.DB.Prepare(`
		DELETE FROM scores WHERE userid = $1;
	`)
	if err != nil {
		return ErrInvalidQueryStatement
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("delete score statement error: %w", err)
	}
	return nil
}

func (record *Record) Put(id string, user string, score string) error {
	stmt1, err := record.DB.Prepare(`
		UPDATE users SET username = $1 WHERE id = $2;
	`)
	if err != nil {
		return ErrInvalidQueryStatement
	}

	stmt2, err := record.DB.Prepare(`
		UPDATE scores SET score = $1 WHERE userid = $2;
	`)
	if err != nil {
		return ErrInvalidQueryStatement
	}

	_, err = stmt1.Exec(user, id)
	if err != nil {
		return fmt.Errorf("update user error: %w", err)
	}
	_, err = stmt2.Exec(score, id)
	if err != nil {
		return fmt.Errorf("update score error: %w", err)
	}
	return nil
}

func (record *Record) Get() ([]Item, error) {
	var username string
	var score string
	var id string
	items := []Item{}
	rows, err := record.DB.Query(`
		SELECT u.id, u.username, s.score FROM users AS u INNER JOIN scores AS s ON u.id = s.userid ORDER BY s.score desc LIMIT 10;
	`)
	if err != nil {
		return nil, ErrInvalidQueryStatement
	}

	for rows.Next() {
		_, err = rows.Scan(&id, &username, &score), nil
		if err != nil {
			return nil, ErrScanCouldntComplete
		}
		item := Item{
			Id:       id,
			Username: username,
			Score:    score,
		}
		items = append(items, item)
	}
	return items, nil
}
