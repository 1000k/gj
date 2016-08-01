package models

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var _db *sql.DB

func getDbHandler(host string) (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", host)
	return
}

func createMessagesTable(db *sql.DB) (res sql.Result, err error) {
	res, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		from_name VARCHAR(255) NOT NULL,
		to_name VARCHAR(255) NOT NULL,
		message TEXT
	)`)
	return
}

// ConnectDb trys to connect specified host and returns connection.
func ConnectDb(host string) (db *sql.DB, err error) {
	db, err = getDbHandler(host)
	if err != nil {
		return
	}
	_, err = createMessagesTable(db)

	_db = db
	return
}

// ResetDb drops all tables and re-creates messages table.
func ResetDb(host string) (res sql.Result, err error) {
	db, err := getDbHandler(host)
	if err != nil {
		return
	}

	res, err = db.Exec(`DROP TABLE messages`)
	if err != nil {
		return
	}

	res, err = createMessagesTable(db)

	return
}

// NewMessage inserts a message record and returns last insert id.
func NewMessage(from, to, message string) (id int64, err error) {
	stmt, err := _db.Prepare("INSERT INTO messages(from_name, to_name, message) values(?,?,?)")
	if err != nil {
		return
	}

	res, err := stmt.Exec(from, to, message)
	if err != nil {
		return
	}

	id, err = res.LastInsertId()
	return
}

type MessageItem struct {
	Id        int64
	FromName  string
	ToName    string
	CreatedAt time.Time
	Message   string
}

type Messages struct {
	Items []MessageItem
}

// FindMessages returns Messages slice.
func FindMessages() (items []MessageItem, err error) {
	rows, err := _db.Query("SELECT id, from_name, to_name, created_at, message FROM messages ORDER BY created_at DESC")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		item := MessageItem{}
		err = rows.Scan(&item.Id, &item.FromName, &item.ToName, &item.CreatedAt, &item.Message)
		if err != nil {
			return
		}
		items = append(items, item)
	}

	return
}

type RankingItem struct {
	Name  string
	Count int
}

type Rankings struct {
	Items []RankingItem
}

func FindRanking(ym string) {

}
