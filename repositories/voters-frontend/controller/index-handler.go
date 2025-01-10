package controller

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"html/template"
	"io/fs"
	"net/http"
)

type FrontendController struct {
	participantRepository domain.ParticipantRepository
	voteRepository        domain.VoteRepository
	context               context.Context
	embedTemplates        fs.FS
}

func NewFrontendController(participantRepository domain.ParticipantRepository, voteRepository domain.VoteRepository, context context.Context, embedTemplates fs.FS) FrontendController {
	return FrontendController{
		participantRepository: participantRepository,
		voteRepository:        voteRepository,
		context:               context,
		embedTemplates:        embedTemplates,
	}
}

func (controller *FrontendController) IndexHandler(responseWriter http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFS(controller.embedTemplates, "index.html")
	if err != nil {
		handleInternalServerError(responseWriter, err)
		return
	}

	// Render the template with the items data
	participants, err1 := controller.participantRepository.FindAll(controller.context)
	if err1 != nil {
		handleInternalServerError(responseWriter, err1)
		return
	}

	err = tmpl.Execute(responseWriter, participants)
	if err != nil {
		handleInternalServerError(responseWriter, err)
		return
	}
}
