package main

import (
	"bbb-voting/prodution-frontend/controller"
	_ "bbb-voting/prodution-frontend/docs"
	postgresqldatamapper "bbb-voting/voting-commons/data-layer/postgresql"
	"bbb-voting/voting-commons/domain"
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
)

//go:embed view/static/*
var staticFilesFull embed.FS

//go:embed view/templates/*
var templatesFull embed.FS

func main() {
	var templates, _ = fs.Sub(templatesFull, "view/templates")
	var staticFiles, _ = fs.Sub(staticFilesFull, "view/static")

	ctx := context.Background()
	postgresqlConnector := postgresqldatamapper.NewPostgresqlConnector(os.Getenv("POSTGRESQL_URI"))
	// redisClient := getRedisClient(os.Getenv("REDIS_URL"))

	var participantRepository domain.ParticipantRepository = postgresqldatamapper.NewParticipantDataMapper(
		postgresqlConnector,
	)
	// var err error
	// participantRepository, err = redisdatamapper.DecorateParticipantDataMapperWithRedis(participantRepository, *redisClient, ctx)
	// if err != nil {
	// 	log.Fatalf("Faled to load cache: %s", err)
	// }

	frontendController := controller.NewFrontendController(
		participantRepository,
		ctx, templates,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", frontendController.GetPage)
	mux.HandleFunc("/votes/totals/thorough", frontendController.GetThoroughTotals)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFiles))))
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	log.Println("Server is running on http://localhost:8081")
	if err := http.ListenAndServe(":8081", corsMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
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

func getRedisClient(url string) *redis.Client {
	opts, err := redis.ParseURL(url)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}
	return redis.NewClient(opts)
}
