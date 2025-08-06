package handlers

import (
	"encoding/json"
	"net/http"
)

type HealthHandler struct {
}

type HealthResponse struct {
	Message string `json:"health"`
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (healthHandler *HealthHandler) GetHealth(wtr http.ResponseWriter, req *http.Request) {
	response := HealthResponse{Message: "Alive"}

	wtr.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(wtr).Encode(response); err != nil {
		http.Error(wtr, err.Error(), http.StatusInternalServerError)
		return
	}
}
