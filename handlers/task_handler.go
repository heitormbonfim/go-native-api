package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/heitormbonfim/go-native-api/models"
)

type TaskHandler struct {
	DB *sql.DB
}

func NewTaskHandler(db *sql.DB) *TaskHandler {
	return &TaskHandler{DB: db}
}

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
	defer rows.Close()

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

	response := TasksResponse{Tasks: tasks}

	wtr.Header().Set("Content-Type", "application/json")
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

	var exists bool
	err = taskHandler.DB.QueryRow(`--sql
	SELECT EXISTS (SELECT 1 FROM tasks WHERE id = $1)
	`, id).Scan(&exists)
	if err != nil {
		http.Error(wtr, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(wtr, "Task not found", http.StatusNotFound)
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
	vars := mux.Vars(req)
	id := vars["id"]

	var exists bool
	err := taskHandler.DB.QueryRow(`--sql
	SELECT EXISTS (SELECT 1 FROM tasks WHERE id = $1)
	`, id).Scan(&exists)
	if err != nil {
		http.Error(wtr, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(wtr, "Task not found", http.StatusNotFound)
		return
	}

	_, err = taskHandler.DB.Exec(`--sql
	DELETE FROM tasks WHERE id = $1
	`, id)
	if err != nil {
		http.Error(wtr, err.Error(), http.StatusInternalServerError)
		return
	}

	wtr.WriteHeader(http.StatusNoContent)
}
