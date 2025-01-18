package controller

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

type FrontendController struct {
	participantRepository domain.ParticipantRepository
	voteRepository        domain.VoteRepository
	ctx                   context.Context
	embedTemplates        fs.FS
	embedStatic           fs.FS
}

func NewFrontendController(participantRepository domain.ParticipantRepository, voteRepository domain.VoteRepository, ctx context.Context, embedTemplates fs.FS, embedStatic fs.FS) FrontendController {
	return FrontendController{
		participantRepository: participantRepository,
		voteRepository:        voteRepository,
		ctx:                   ctx,
		embedTemplates:        embedTemplates,
		embedStatic:           embedStatic,
	}
}

func (frontendController *FrontendController) GetServerMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", frontendController.IndexHandler)
	mux.HandleFunc("/pages/totals/rough", frontendController.LoadRoughTotalPage)

	mux.HandleFunc("/votes", frontendController.PostVoteHandler)
	mux.HandleFunc("/participants", frontendController.GetParticipantsHandler)
	mux.HandleFunc("/votes/totals/rough", frontendController.GetVotesRoughTotalsHandler)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(frontendController.embedStatic))))
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return mux
}



func loadBody(responseWriter http.ResponseWriter, request *http.Request, contentBody any) {
	bytesBody, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(responseWriter, "Error when read", http.StatusBadRequest)
	}

	err = json.Unmarshal(bytesBody, &contentBody)
	if err != nil {
		http.Error(responseWriter, "Error when read", http.StatusMethodNotAllowed)
	}
}

func handleInternalServerError(responseWriter http.ResponseWriter, err error) {
	http.Error(responseWriter, "Internal Server Error", http.StatusInternalServerError)
	log.Fatal(err)
}
