package controller

import (
	usercases "bbb-voting/voting-commons/user-cases"
	"io/fs"
	"net/http"

	"bbb-voting/prodution-frontend/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

type FrontendController struct {
	getThoroughTotalsUserCase usercases.GetThoroughTotalsUserCase
	templates                 fs.FS
	staticFiles               fs.FS
}

func NewFrontendController(getThoroughTotalsUserCase usercases.GetThoroughTotalsUserCase, templates fs.FS, staticFiles fs.FS) FrontendController {
	return FrontendController{
		getThoroughTotalsUserCase: getThoroughTotalsUserCase,
		templates:                 templates,
		staticFiles:               staticFiles,
	}
}

func (frontendController *FrontendController) GetServerMux() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", frontendController.GetPage)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(frontendController.staticFiles))))

	mux.HandleFunc("/api/votes/totals/thorough", frontendController.GetThoroughTotals)

	configSwagger()
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return corsMiddleware(mux)
}

func configSwagger() {
	docs.SwaggerInfo.Title = "Prodution Frontend"
	docs.SwaggerInfo.Description = "Frontend for Prodution"
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
