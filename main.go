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

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS gjs (
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

		initializeDb()

		db, err := sql.Open("sqlite3", "./gj.db")
		checkErr(err, "Failed to open database file.")

		stmt, err := db.Prepare("INSERT INTO gjs(from_name, to_name, message) values(?,?,?)")
		checkErr(err, "Prepared query is invalid.")

		res, err := stmt.Exec(r.FormValue("from_name"), r.FormValue("to_name"), r.FormValue("message"))
		checkErr(err, "Cannot execute query.")

		id, err := res.LastInsertId()
		checkErr(err, "Cannot fetch last id.")

		log.Println("id: ", id)
		tv.Message = "Saved"
	}

	t := template.Must(template.ParseFiles("templates/index.html"))
	t.Execute(w, tv)
}

func main() {
	http.HandleFunc("/", IndexHandler)
	err := http.ListenAndServe(":8000", nil)
	checkErr(err, "ListenAndServer: ")
}
