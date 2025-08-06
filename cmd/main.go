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
	// endpoints
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Alive!"))
	}).Methods("GET")

	// Tasks Endpoints
	taskHandler := handlers.NewTaskHandler(dbConnection)
	router.HandleFunc("/tasks", taskHandler.GetTasks).Methods("GET")

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
		return
	}
}
