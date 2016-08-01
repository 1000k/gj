package models

import (
	"database/sql"
	"fmt"
	"regexp"
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

type NewMessageParams struct {
	From, To, Message, CreatedAt string
}

// NewMessage inserts a message record and returns last insert id.
func NewMessage(p NewMessageParams) (id int64, err error) {
	stmt, err := _db.Prepare("INSERT INTO messages(from_name, to_name, message, created_at) values(?,?,?,?)")
	if err != nil {
		return
	}

	if p.CreatedAt == "" {
		p.CreatedAt = time.Now().Format("2016-01-02 03:04:05")
	}

	res, err := stmt.Exec(p.From, p.To, p.Message, p.CreatedAt)
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

func getToDate(from string) (to string) {
	t, _ := time.Parse("2006-01-02", from)
	to = t.AddDate(0, 1, 0).Format("2006-01-02")
	return
}

// FindRanking returns ranking slice.
// Parameter 'ym' should be as `201608`
func FindRanking(ym string) (items []RankingItem, err error) {
	rep := regexp.MustCompile(`(\d{4})(\d{2})`)
	fromDate := rep.ReplaceAllString(ym, "$1-$2-01")
	if fromDate == "" {
		err = fmt.Errorf("Parameter 'ym' should be as 'YYYYMM'")
		return
	}
	toDate := getToDate(fromDate)

	stmt, err := _db.Prepare("SELECT to_name, COUNT(*) AS cnt FROM messages WHERE created_at >= ? AND created_at < ? GROUP BY to_name ORDER BY cnt DESC LIMIT 10")
	if err != nil {
		return
	}

	rows, err := stmt.Query(fromDate, toDate)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		item := RankingItem{}
		err = rows.Scan(&item.Name, &item.Count)
		if err != nil {
			return
		}
		items = append(items, item)
	}

	return
}
