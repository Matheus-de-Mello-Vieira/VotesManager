package controller

import (
	"bbb-voting/voting-commons/domain"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (controller *FrontendController) VoteCastingHandler(responseWriter http.ResponseWriter, request *http.Request) {
	// Only allow POST requests
	if request.Method != http.MethodPost {
		http.Error(responseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data from POST request
	if err := request.ParseForm(); err != nil {
		http.Error(responseWriter, "Unable to parse form", http.StatusBadRequest)
		log.Fatal(err)
		return
	}

	// Get the participant ID from the form
	idStr := request.FormValue("id")
	if idStr == "" {
		http.Error(responseWriter, "Missing participant ID", http.StatusBadRequest)
		return
	}

	// Convert ID to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(responseWriter, "Invalid participant ID", http.StatusBadRequest)
		return
	}

	participant, _ := controller.participantRepository.FindByID(controller.context, id)

	if participant == nil {
		http.Error(responseWriter, "Participant not found", http.StatusNotFound)
		return
	}

	vote := domain.Vote{Participant: *participant, Timestamp: time.Now()}

	controller.voteRepository.SaveOne(controller.context, &vote)

	controller.loadRoughTotalPage(responseWriter)
}

type RoughTotalPresenter struct {
	Labels []string
	Votes  []int
}

func loadRoughTotalPresenter(totalsMap map[domain.Participant]int) RoughTotalPresenter {
	result := RoughTotalPresenter{
		Labels: []string{},
		Votes:  []int{},
	}

	for participant, vote := range totalsMap {
		result.Labels = append(result.Labels, participant.Name)
		result.Votes = append(result.Votes, vote)
	}

	return result
}

func (controller *FrontendController) loadRoughTotalPage(responseWriter http.ResponseWriter) {
	tmpl, err := template.ParseFS(controller.embedTemplates, "rough_results.html")
	if err != nil {
		handleInternalServerError(responseWriter, err)
		return
	}

	totalsMap, err1 := controller.participantRepository.GetRoughTotals(controller.context)
	if err1 != nil {
		handleInternalServerError(responseWriter, err)
		return
	}

	data := loadRoughTotalPresenter(totalsMap)

	err = tmpl.Execute(responseWriter, data)
	if err != nil {
		handleInternalServerError(responseWriter, err)
	}
}
