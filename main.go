package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

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

		params := models.NewMessageParams{From: r.FormValue("from_name"), To: r.FormValue("to_name"), Message: r.FormValue("message")}
		id, err := models.NewMessage(params)
		checkErr(err)

		log.Printf("New message saved. id: %v, values: %v\n", id, r.Form)
		tv.Message = "Saved"
	}

	t := template.Must(template.ParseFiles("templates/index.html", "templates/_header.html"))
	t.Execute(w, tv)
}

func ChartHandler(w http.ResponseWriter, r *http.Request) {
	_, err := models.ConnectDb(dbfile)
	checkErr(err)

	thisMonth := time.Now().Format("200601")
	ranking, err := models.FindRanking(thisMonth)
	checkErr(err)

	labels := []string{}
	data := []int{}
	for _, v := range ranking {
		labels = append(labels, v.Name)
		data = append(data, v.Count)
	}

	l, _ := json.Marshal(labels)
	d, _ := json.Marshal(data)
	out := map[string]interface{}{
		"Labels": template.JS(string(l)),
		"Data":   template.JS(string(d)),
	}
	// log.Println(labels, data, out)

	t := template.Must(template.ParseFiles("templates/chart.html", "templates/_header.html"))
	t.Execute(w, out)
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

	messages, err := models.FindMessages()
	checkErr(err)

	t := template.Must(template.ParseFiles("templates/messages.html", "templates/_header.html"))
	err = t.Execute(w, messages)
	checkErr(err)
}

func main() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/chart", ChartHandler)
	http.HandleFunc("/messages", MessagesHandler)
	err := http.ListenAndServe(":8000", nil)
	checkErr(err, "ListenAndServer: ")
}
