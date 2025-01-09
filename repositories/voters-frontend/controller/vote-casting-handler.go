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

	if !verifyRecaptcha() {
		http.Error(responseWriter, "CAPTCHA verification failed. Please try again.", http.StatusForbidden)
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

func verifyRecaptcha() bool {
	return true
}

func (controller *FrontendController) loadRoughTotalPage(responseWriter http.ResponseWriter) {
	tmpl, err := template.ParseFiles(templatesPath + "rough_results.html")
	if err != nil {
		handleInternalServerError(responseWriter, err)
		return
	}

	totalsMap, err1 := controller.participantRepository.GetRoughTotals(controller.context)
	if err1 != nil {
		handleInternalServerError(responseWriter, err)
		return
	}

	data := struct {
		Labels []string
		Votes  []int
	}{
		Labels: []string{},
		Votes:  []int{},
	}

	for participant, vote := range totalsMap {
		data.Labels = append(data.Labels, participant.Name)
		data.Votes = append(data.Votes, vote)
	}

	err = tmpl.Execute(responseWriter, data)
	if err != nil {
		handleInternalServerError(responseWriter, err)
	}
}
