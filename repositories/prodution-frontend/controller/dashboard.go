package controller

import (
	"bbb-voting/voting-commons/domain"
	"context"
	"encoding/json"
	"net/http"
	"log"
)

type FrontendController struct {
	participantRepository domain.ParticipantRepository
	context context.Context
}

func NewFrontendController(participantRepository domain.ParticipantRepository, context context.Context) FrontendController {
	return FrontendController{
		participantRepository: participantRepository,
		context: context,
	}
}

type ThoroughTotalsResponseModel struct {
	GeneralTotal       int 
	TotalByHour        []domain.TotalByHour
	TotalByParticipant map[string]int
}


func (controller *FrontendController) GetThoroughTotals (responseWriter http.ResponseWriter, request *http.Request) {
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
		GeneralTotal: content.GeneralTotal,
		TotalByHour: content.TotalByHour,
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
