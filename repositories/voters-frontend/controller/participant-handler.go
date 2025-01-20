package controller

import (
	"encoding/json"
	"html/template"
	"net/http"
)

// @Summary Serve HTML index page
// @Description Responds with an HTML page with the index page
// @Tags html
// @Produce html
// @Success 200 {string} string "HTML Content"
// @Router / [get]
func (controller *FrontendController) IndexHandler(responseWriter http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFS(controller.embedTemplates, "index.html")
	if err != nil {
		handleInternalServerError(responseWriter, err)
		return
	}

	err = tmpl.Execute(responseWriter, nil)
	if err != nil {
		handleInternalServerError(responseWriter, err)
		return
	}
}

// @Summary Get Participants
// @Description Responds with the list of participants
// @Tags participants
// @Accept json
// @Produce json
// @Success 200 {object} []domain.Participant
// @Router /participants [get]
func (controller *FrontendController) GetParticipantsHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(responseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	participants, err := controller.getParticipantsUserCase.Execute()
	if err != nil {
		handleInternalServerError(responseWriter, err)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(responseWriter).Encode(participants)
}
