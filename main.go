package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type TemplateValues struct {
	HasError bool
	Message  string
}

func checkErr(err error, messages ...string) {
	if err != nil {
		tv.HasError = false
		tv.Message = strings.Join(messages, " ")
		panic(err)
	}
}

func initializeDb() *sql.DB {
	db, err := sql.Open("sqlite3", "./gj.db")
	checkErr(err)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		from_name VARCHAR(255) NOT NULL,
		to_name VARCHAR(255) NOT NULL,
		message TEXT
	)`)
	checkErr(err)

	return db
}

var tv = TemplateValues{
	HasError: true,
	Message:  "",
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method == "POST" {
		r.ParseForm()
		log.Println(r.Form)

		db := initializeDb()

		stmt, err := db.Prepare("INSERT INTO messages(from_name, to_name, message) values(?,?,?)")
		checkErr(err, "Prepared query is invalid.")

		res, err := stmt.Exec(r.FormValue("from_name"), r.FormValue("to_name"), r.FormValue("message"))
		checkErr(err, "Cannot execute query.")

		id, err := res.LastInsertId()
		checkErr(err, "Cannot fetch last id.")

		log.Println("id: ", id)
		tv.Message = "Saved"
	}

	t := template.Must(template.ParseFiles("templates/index.html", "templates/_header.html"))
	t.Execute(w, tv)
}

func ChartHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/chart.html", "templates/_header.html"))
	t.Execute(w, nil)
}

type MessageItem struct {
	Id        int
	FromName  string
	ToName    string
	CreatedAt string
	Message   string
}

func MessagesHandler(w http.ResponseWriter, r *http.Request) {
	db := initializeDb()

	rows, err := db.Query("SELECT id, from_name, to_name, created_at, message FROM messages ORDER BY created_at DESC")
	checkErr(err, "Prepared query is invalid.")
	defer rows.Close()

	var result []MessageItem

	for rows.Next() {
		item := MessageItem{}
		err = rows.Scan(&item.Id, &item.FromName, &item.ToName, &item.CreatedAt, &item.Message)
		checkErr(err)
		result = append(result, item)
	}

	t := template.Must(template.ParseFiles("templates/messages.html", "templates/_header.html"))
	err = t.Execute(w, result)
	checkErr(err)
}

func main() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/chart", ChartHandler)
	http.HandleFunc("/messages", MessagesHandler)
	err := http.ListenAndServe(":8000", nil)
	checkErr(err, "ListenAndServer: ")
}
