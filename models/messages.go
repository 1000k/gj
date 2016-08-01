package models

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var _db *sql.DB

func checkErr(err error, messages ...string) {
	if err != nil {
		panic(err)
	}
}

// ConnectDb trys to connect specified host and returns connection.
func ConnectDb(host string) (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", host)
	if err != nil {
		return db, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		from_name VARCHAR(255) NOT NULL,
		to_name VARCHAR(255) NOT NULL,
		message TEXT
	)`)

	_db = db
	return db, err
}

// NewMessage inserts a message record and returns last insert id.
func NewMessage(from, to, message string) (id int64, err error) {
	stmt, err := _db.Prepare("INSERT INTO messages(from_name, to_name, message) values(?,?,?)")
	if err != nil {
		return id, err
	}

	res, err := stmt.Exec(from, to, message)
	if err != nil {
		return id, err
	}

	id, err = res.LastInsertId()
	return id, err
}

type MessageItem struct {
	Id        int
	FromName  string
	ToName    string
	CreatedAt string
	Message   string
}

type Messages struct {
	Items []MessageItem
}

// FindMessages returns Messages slice.
func FindMessages() (items []MessageItem, err error) {
	item := MessageItem{Id: 1, FromName: "Moja", ToName: "Uhouho", Message: "HEY YO", CreatedAt: "2016-08-01"}
	items = append(items, item)
	return items, fmt.Errorf("mojamoja")
}
