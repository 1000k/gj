package main

import (
	"html/template"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	t := template.Must(template.ParseFiles("templates/index.html"))
	t.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", IndexHandler)
	http.ListenAndServe(":8000", nil)
}
