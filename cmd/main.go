package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/heitormbonfim/go-native-api/config"
	"github.com/heitormbonfim/go-native-api/handlers"
	"github.com/heitormbonfim/go-native-api/models"
)

func main() {
	dbConnection := config.SetupDB()
	defer dbConnection.Close()

	_, err := dbConnection.Exec(models.CreateTableSQL)
	if err != nil {
		log.Fatal(err)
		return
	}

	router := mux.NewRouter()
	// Endpoints

	// Health
	healthHandler := handlers.NewHealthHandler()
	router.HandleFunc("/health", healthHandler.GetHealth).Methods("GET")

	// Tasks Endpoints
	taskHandler := handlers.NewTaskHandler(dbConnection)
	router.HandleFunc("/tasks", taskHandler.GetTasks).Methods("GET")
	router.HandleFunc("/task", taskHandler.CreateTask).Methods("POST")
	router.HandleFunc("/task/{id}", taskHandler.UpdateTask).Methods("PUT")
	router.HandleFunc("/task/{id}", taskHandler.DeleteTask).Methods("DELETE")

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
		return
	}
}
