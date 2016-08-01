package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/1000k/gj/models"
	_ "github.com/mattn/go-sqlite3"
)

const dbfile = "./gj.db"

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

var tv = TemplateValues{
	HasError: true,
	Message:  "",
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method == "POST" {
		r.ParseForm()

		_, err := models.ConnectDb(dbfile)
		checkErr(err)

		id, err := models.NewMessage(r.FormValue("from_name"), r.FormValue("to_name"), r.FormValue("message"))
		checkErr(err)

		log.Printf("New message saved. id: %v, values: %v\n", id, r.Form)
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
	_, err := models.ConnectDb(dbfile)
	checkErr(err)

	result, err := models.FindMessages()
	checkErr(err)

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
