package controller

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"encoding/json"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

type FrontendController struct {
	participantRepository domain.ParticipantRepository
	context               context.Context
	embedTemplates        fs.FS
}

func NewFrontendController(participantRepository domain.ParticipantRepository, context context.Context, embedTemplates fs.FS) FrontendController {
	return FrontendController{
		participantRepository: participantRepository,
		context:               context,
		embedTemplates:        embedTemplates,
	}
}

type ThoroughTotalsResponseModel struct {
	GeneralTotal       int                  `json:"general_total"`
	TotalByHour        []domain.TotalByHour `json:"total_by_hour"`
	TotalByParticipant map[string]int       `json:"total_by_participant"`
}

// @Summary Serve HTML thorough total page
// @Description Responds with an HTML page with a thorough total graph
// @Tags html
// @Produce html
// @Success 200 {string} string "HTML Content"
// @Router / [get]
func (controller *FrontendController) GetPage(responseWriter http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFS(controller.embedTemplates, "dashboard.html")

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

func handleInternalServerError(responseWriter http.ResponseWriter, err error) {
	http.Error(responseWriter, "Internal Server Error", http.StatusInternalServerError)
	log.Fatal(err)
}

// @Summary Get Thorough Totals
// @Description Get throught totals
// @Tags totals votes
// @Accept  json
// @Produce  json
// @Success 200 {object} ThoroughTotalsResponseModel
// @Router /votes/totals/thorough [get]
func (controller *FrontendController) GetThoroughTotals(responseWriter http.ResponseWriter, request *http.Request) {
	// Only allow GET requests
	if request.Method != http.MethodGet {
		http.Error(responseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	content, err := controller.participantRepository.GetThoroughTotals(controller.context)
	if err != nil {
		log.Printf("Error on get result: %v", err)
		return
	}

	responseModel := ThoroughTotalsResponseModel{
		GeneralTotal:       content.GeneralTotal,
		TotalByHour:        content.TotalByHour,
		TotalByParticipant: map[string]int{},
	}

	for participant, value := range content.TotalByParticipant {
		responseModel.TotalByParticipant[participant.Name] = value
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(responseWriter).Encode(responseModel)
}
