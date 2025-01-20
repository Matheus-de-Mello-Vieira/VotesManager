package controller

import (
	usercases "bbb-voting/voting-commons/user-cases"
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
	mux.HandleFunc("/pages/totals/rough", frontendController.LoadRoughTotalPage)

	mux.HandleFunc("/votes", frontendController.PostVoteHandler)
	mux.HandleFunc("/participants", frontendController.GetParticipantsHandler)
	mux.HandleFunc("/votes/totals/rough", frontendController.GetVotesRoughTotalsHandler)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(frontendController.embedStatic))))
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return corsMiddleware(mux)
}

func configSwagger() {
	docs.SwaggerInfo.Title = "Voters Frontend"
	docs.SwaggerInfo.Description = "Frontend for Voters"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
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

func loadBody(responseWriter http.ResponseWriter, request *http.Request, contentBody any) {
	bytesBody, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(responseWriter, "Error when read", http.StatusBadRequest)
	}

	err = json.Unmarshal(bytesBody, &contentBody)
	if err != nil {
		http.Error(responseWriter, "Error when read", http.StatusMethodNotAllowed)
	}
}

func handleInternalServerError(responseWriter http.ResponseWriter, err error) {
	http.Error(responseWriter, "Internal Server Error", http.StatusInternalServerError)
	log.Fatal(err)
}
