package controller

import (
	"bbb-voting/voting-commons/domain"
	"encoding/json"
	"html/template"
	"net/http"
	"time"
)

func (controller *FrontendController) GetVotesRoughTotalsHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(responseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	totalsMap, err := controller.participantRepository.GetRoughTotals(controller.context)
	if err != nil {
		handleInternalServerError(responseWriter, err)
		return
	}

	result := formatRoughTotals(totalsMap)

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(responseWriter).Encode(result)
}
func formatRoughTotals(totalsMap map[domain.Participant]int) map[string]int {
	result := map[string]int{}

	for participant, vote := range totalsMap {
		result[participant.Name] = vote
	}

	return result
}

type postVoteBodyModel struct {
	ParticipantID int `json:"participant_id"`
}

func (controller *FrontendController) PostVoteHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(responseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body := postVoteBodyModel{}
	loadBody(responseWriter, request, &body)

	participant, _ := controller.participantRepository.FindByID(controller.context, body.ParticipantID)

	if participant == nil {
		http.Error(responseWriter, "Participant not found", http.StatusNotFound)
		return
	}

	vote := domain.Vote{Participant: *participant, Timestamp: time.Now()}

	controller.voteRepository.SaveOne(controller.context, &vote)

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusCreated)
	json.NewEncoder(responseWriter).Encode(vote)
}

func (controller *FrontendController) LoadRoughTotalPage(responseWriter http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFS(controller.embedTemplates, "rough_results.html")
	if err != nil {
		handleInternalServerError(responseWriter, err)
		return
	}

	err = tmpl.Execute(responseWriter, nil)
	if err != nil {
		handleInternalServerError(responseWriter, err)
	}
}
