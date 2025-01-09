package controller

import (
	"html/template"
	"log"
	"net/http"
)

const templatesPath = "voters-frontend/view/templates/"

type Item struct {
	Name string
}

var items = []Item{
	{"Item 1"},
	{"Item 2"},
	{"Item 3"},
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(templatesPath + "index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	// Render the template with the items data
	err = tmpl.Execute(w, items)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal(err)
	}
}
