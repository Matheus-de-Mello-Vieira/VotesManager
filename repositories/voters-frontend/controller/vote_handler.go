package controller

import (
	"bbb-voting/voting-commons/domain"
	usercases "bbb-voting/voting-commons/user-cases"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

// @Summary Get Rough Totals
// @Description Get rough totals
// @Tags api
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]int
// @Router /api/votes/totals/rough [get]
func (controller *FrontendController) GetVotesRoughTotalsHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(responseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	totalsMap, err := controller.getRoughTotalsUserCase.Execute()
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

// @Summary Post Vote
// @Description Cast a Vote
// @Tags api
// @Accept  json
// @Produce  json
// @Body postVoteBodyModel
// @Success 201 {object} domain.Vote
// @Router /api/votes [post]
func (controller *FrontendController) PostVoteHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(responseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body := usercases.CastVoteDTO{}
	err := loadBody(responseWriter, request, &body)
	if err != nil {
		return
	}

	vote, err := controller.castVoteUserCase.Execute(&body)
	if err != nil {
		if errors.Is(err, usercases.ErrParticipantNotFound) {
			http.Error(responseWriter, fmt.Sprint(err), http.StatusNotFound)
		} else {
			handleInternalServerError(responseWriter, err)
		}
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusCreated)
	json.NewEncoder(responseWriter).Encode(vote)
}

// @Summary Serve HTML rought total page
// @Description Responds with an HTML page with a rought total graph
// @Tags html
// @Produce html
// @Success 200 {string} string "HTML Content"
// @Router /after-vote [get]
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
