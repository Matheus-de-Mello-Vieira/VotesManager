package controller

import (
	"bbb-voting/voting-commons/domain"
	"html/template"
	"net/http"
	"strconv"
)


func VoteCastingHandler(w http.ResponseWriter, r *http.Request) {
    // Only allow POST requests
    if r.Method != http.MethodPost {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse form data from POST request
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

    // Get the participant ID from the form
    idStr := r.FormValue("id")
    if idStr == "" {
        http.Error(w, "Missing participant ID", http.StatusBadRequest)
        return
    }

    // Convert ID to integer
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid participant ID", http.StatusBadRequest)
        return
    }

    // Find the participant
    var participant *domain.Participant
    for _, p := range participants {
        if p.ParticipantID == id {
            participant = &p
            break
        }
    }

    if participant == nil {
        http.Error(w, "Participant not found", http.StatusNotFound)
        return
    }

    // Load and execute template for single participant
    tmpl, err := template.ParseFiles("templates/item.html")
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    err = tmpl.Execute(w, participant)
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}