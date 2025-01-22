package controller

import (
	usercases "bbb-voting/voting-commons/user-cases"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"net/http"

	"bbb-voting/voters-frontend/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

type FrontendController struct {
	getRoughTotalsUserCase  usercases.GetRoughTotalsUserCase
	getParticipantsUserCase usercases.GetParticipantsUserCase
	castVoteUserCase        usercases.CastVoteUserCase
	embedTemplates          fs.FS
	embedStatic             fs.FS
}

func NewFrontendController(getRoughTotalsUserCase usercases.GetRoughTotalsUserCase, getParticipantsUserCase usercases.GetParticipantsUserCase, castVoteUserCase usercases.CastVoteUserCase, ctx context.Context, embedTemplates fs.FS, embedStatic fs.FS) FrontendController {
	return FrontendController{
		getRoughTotalsUserCase:  getRoughTotalsUserCase,
		getParticipantsUserCase: getParticipantsUserCase,
		castVoteUserCase:        castVoteUserCase,
		embedTemplates:          embedTemplates,
		embedStatic:             embedStatic,
	}
}

func (frontendController *FrontendController) GetServerMux() http.Handler {
	configSwagger()

	mux := http.NewServeMux()
	mux.HandleFunc("/", frontendController.IndexHandler)
	mux.HandleFunc("/after-vote", frontendController.LoadRoughTotalPage)

	mux.HandleFunc("/api/votes", captchaMiddleware(frontendController.PostVoteHandler))
	mux.HandleFunc("/api/participants", frontendController.GetParticipantsHandler)
	mux.HandleFunc("/api/votes/totals/rough", frontendController.GetVotesRoughTotalsHandler)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(frontendController.embedStatic))))
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return corsMiddleware(mux)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins; restrict as needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func configSwagger() {
	docs.SwaggerInfo.Title = "Voters Frontend"
	docs.SwaggerInfo.Description = "Frontend for Voters"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}

func loadBody(responseWriter http.ResponseWriter, request *http.Request, contentBody any) error {
	bytesBody, err := io.ReadAll(request.Body)
	if err != nil {
		log.Printf("Error when read the body: %s", err)
		http.Error(responseWriter, "Error when read the body", http.StatusBadRequest)
		return err
	}

	request.Body = io.NopCloser(bytes.NewBuffer(bytesBody))

	err = json.Unmarshal(bytesBody, &contentBody)
	if err != nil {
		log.Printf("Error when read the body: %s", err)
		http.Error(responseWriter, "Error when read the body", http.StatusMethodNotAllowed)
		return err
	}

	return err
}

func handleInternalServerError(responseWriter http.ResponseWriter, err error) {
	http.Error(responseWriter, "Internal Server Error", http.StatusInternalServerError)
	log.Println(err)
}
