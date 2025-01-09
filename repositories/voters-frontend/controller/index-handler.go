package controller

import (
	"html/template"
	"log"
	"net/http"
	"bbb-voting/voting-commons/domain"
)

const templatesPath = "voters-frontend/view/templates/"

var participants = []domain.Participant{
	{ParticipantID: 1, Name: "Isaac Newton"},
	{ParticipantID: 2, Name: "Albert Einstein"},
	{ParticipantID: 3, Name: "Marie Curie"},
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(templatesPath + "index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	// Render the template with the items data
	err = tmpl.Execute(w, participants)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal(err)
	}
}
