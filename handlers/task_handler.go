package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/heitormbonfim/go-native-api/models"
)

// The original TaskHandler and its methods...
type TaskHandler struct {
	DB *sql.DB
}

// Handler "Constructor"
func NewTaskHandler(db *sql.DB) *TaskHandler {
	return &TaskHandler{DB: db}
}

// Handlers "Methods"

// Response struct to match the desired JSON output
type TasksResponse struct {
	Tasks []models.Task `json:"tasks"`
}

func (taskHandler *TaskHandler) GetTasks(wtr http.ResponseWriter, req *http.Request) {
	rows, err := taskHandler.DB.Query(`--sql
  SELECT * FROM tasks;
  `)
	if err != nil {
		http.Error(wtr, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close() // Ensure rows are closed

	// getting all the rows from the result of the query and adding them into an array
	var tasks []models.Task

	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status)
		if err != nil {
			http.Error(wtr, err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	// Create an instance of the new struct and assign the tasks slice
	response := TasksResponse{Tasks: tasks}

	// sending ok response implicitly
	wtr.Header().Set("Content-Type", "application/json")
	// Encode the new struct instead of the slice directly
	// the statement below "trycatch" for the enconder
	if err := json.NewEncoder(wtr).Encode(response); err != nil {
		http.Error(wtr, err.Error(), http.StatusInternalServerError)
	}
}

func (taskHandler *TaskHandler) CreateTask(wtr http.ResponseWriter, req *http.Request) {
	var task models.Task
	err := json.NewDecoder(req.Body).Decode(&task)
	if err != nil {
		http.Error(wtr, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = taskHandler.DB.Exec(`--sql
	INSERT INTO tasks (title, description, status)
	VALUES ($1, $2, $3)
	`, &task.Title, &task.Description, &task.Status)
	if err != nil {
		http.Error(wtr, err.Error(), http.StatusInternalServerError)
		return
	}

	wtr.Header().Set("Content-Type", "application/json")
	wtr.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(wtr).Encode(task); err != nil {
		http.Error(wtr, err.Error(), http.StatusInternalServerError)
	}
}

func (taskHandler *TaskHandler) UpdateTask(wtr http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	var task models.Task

	err := json.NewDecoder(req.Body).Decode(&task)
	if err != nil {
		http.Error(wtr, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = taskHandler.DB.Exec(`--sql
	UPDATE tasks 
	SET title = $1, description = $2, status = $3
	WHERE id = $4 
	`, &task.Title, &task.Description, &task.Status, id)
	if err != nil {
		http.Error(wtr, err.Error(), http.StatusInternalServerError)
		return
	}

	wtr.Header().Set("Content-Type", "application/json")
	wtr.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(wtr).Encode(task); err != nil {
		http.Error(wtr, err.Error(), http.StatusInternalServerError)
	}
}

func (taskHandler *TaskHandler) DeleteTask(wtr http.ResponseWriter, req *http.Request) {

}
