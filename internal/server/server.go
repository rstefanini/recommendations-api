package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"

	"recommendation/internal/database"
	"recommendation/internal/handler"
	"recommendation/internal/repository"
	"recommendation/internal/service"
)

type Server struct {
	Database                     Database
	UserInteractionsHandler      UserInteractionsHandler
	ProductRecommendationHandler ProductRecommendationHandler
}

type Database interface {
	Health() map[string]string
	Client() *mongo.Client
	Database() *mongo.Database
	Collection() *mongo.Collection
}

type UserInteractionsHandler interface {
	Collector(w http.ResponseWriter, r *http.Request)
}

type ProductRecommendationHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	newServer := buildServer()

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func buildServer() *Server {
	db := database.New()

	userBehaviorRepository := repository.NewUserBehaviorRepository(db.Collection())
	userInteractionService := service.NewUserInteractionService(userBehaviorRepository)
	userInteractionHandler := handler.NewUserInteractionHandler(userInteractionService)

	productRecommendationService := service.NewProductRecommendationService(userBehaviorRepository)
	productRecommendationHandler := handler.NewProductRecommendationHandler(productRecommendationService)

	return &Server{
		Database:                     db,
		UserInteractionsHandler:      userInteractionHandler,
		ProductRecommendationHandler: productRecommendationHandler,
	}
}

func (s *Server) logHttp(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("Started %s %s", r.Method, r.URL.Path)
		handler.ServeHTTP(w, r)
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	}
}

func (s *Server) asJSON(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		handler.ServeHTTP(w, r)
	}
}
