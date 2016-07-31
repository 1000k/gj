package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	// "github.com/kr/pretty"
)

func checkErr(err error, messages ...string) {
	if err != nil {
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

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method == "POST" {
		r.ParseForm()
		log.Println(r.Form)

		initializeDb()

		db, err := sql.Open("sqlite3", "./gj.db")
		checkErr(err)

		stmt, err := db.Prepare("INSERT INTO gjs(from_name, to_name, message) values(?,?,?)")
		checkErr(err)

		res, err := stmt.Exec(r.FormValue("from_name"), r.FormValue("to_name"), r.FormValue("message"))
		checkErr(err)

		id, err := res.LastInsertId()
		checkErr(err)

		log.Println("id: ", id)
	}

	t := template.Must(template.ParseFiles("templates/index.html"))
	t.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", IndexHandler)
	err := http.ListenAndServe(":8000", nil)
	checkErr(err, "ListenAndServer: ")
}
