package controller

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type FrontendController struct {
	participantRepository domain.ParticipantRepository
	context               context.Context
}

func NewFrontendController(participantRepository domain.ParticipantRepository, context context.Context) FrontendController {
	return FrontendController{
		participantRepository: participantRepository,
		context:               context,
	}
}

const templatesPath = "prodution-frontend/view/templates/"

type ThoroughTotalsResponseModel struct {
	GeneralTotal       int
	TotalByHour        []domain.TotalByHour
	TotalByParticipant map[string]int
}

func (controller *FrontendController) GetPage(responseWriter http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFiles(templatesPath + "dashboard.html")
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
	result, err1 := json.Marshal(responseModel)
	if err1 != nil {
		log.Printf("Error marshalling vote: %v", err1)
		return
	}

	responseWriter.Write(result)
}
