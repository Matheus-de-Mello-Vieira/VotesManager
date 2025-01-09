package controller

import (
	"html/template"
	"log"
	"net/http"
	"bbb-voting/voting-commons/domain"
	"context"
)

const templatesPath = "voters-frontend/view/templates/"

type FrontendController struct {
	participantRepository domain.ParticipantRepository
	voteRepository domain.VoteRepository
	context context.Context
}

func NewFrontendController(participantRepository domain.ParticipantRepository, voteRepository domain.VoteRepository, context context.Context) FrontendController {
	return FrontendController{
		participantRepository: participantRepository,
		voteRepository: voteRepository,
		context: context,
	}
}

func (controller *FrontendController) IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(templatesPath + "index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	// Render the template with the items data
	participants, err1 := controller.participantRepository.FindAll(controller.context)
	if err1 != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal(err1)
		return
	}

	err = tmpl.Execute(w, participants)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal(err)
	}
}
