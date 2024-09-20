package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("/", s.notFoundHandler)

	mux.HandleFunc("/health", s.asJSON(s.healthHandler))

	mux.HandleFunc("POST /collector/interaction", s.logHttp(s.asJSON(s.UserInteractionsHandler.Collector)))

	mux.HandleFunc("GET /recommendations/users/{user_id}", s.logHttp(s.asJSON(s.ProductRecommendationHandler.Get)))

	return mux
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["health"] = "Healthy"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
