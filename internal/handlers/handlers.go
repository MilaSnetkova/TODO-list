package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MilaSnetkova/TODO-list/internal/models"
	"github.com/MilaSnetkova/TODO-list/internal/service"

)

type TaskHandler struct {
	TaskService service.TaskService
}

func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{
		TaskService: taskService,
	}
}

// Обработчик для маршрута /api/task
func (h *TaskHandler) TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		id := r.URL.Query().Get("id")
		if id == "" {
			// Если id нет, выполняем добавление задачи
			h.AddTaskHandler(w, r)
		} else {
			// Если id есть, отмечаем задачу как выполненную
			h.DoneTaskHandler(w, r)
		}
	case http.MethodGet:
		h.GetTaskInfoHandler(w, r)
	case http.MethodPut:
		h.UpdateTaskHandler(w, r)
	case http.MethodDelete:
		h.DeleteTaskHandler(w, r)
	default:
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

// Добавление задачи 
func (h *TaskHandler) AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, `{"error":"Failed to decode JSON"}`, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, `{"error":"Empty title field"}`, http.StatusBadRequest)
		return
	}

	id, err := h.TaskService.AddTask(&task)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
}

// Получение списка задач с фильтром
func (h *TaskHandler) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	id := r.URL.Query().Get("id") // получаем id из параметров запроса

	tasks, err := h.TaskService.GetTasks(search, id) // передаем search и id
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"tasks": tasks})
}

// Получение информации о задаче
func (h *TaskHandler) GetTaskInfoHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"id is required"}`, http.StatusBadRequest)
		return
	}

	tasks, err := h.TaskService.GetTasks("", id)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	if len(tasks) == 0 {
		http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
		return
	}

	task := tasks[0]

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, `{"error":"failed to encode task"}`, http.StatusInternalServerError)
	}
}

// Обновление задачи
func (h *TaskHandler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, `{"error":"Failed to decode JSON"}`, http.StatusBadRequest)
		return
	}

	if err := h.TaskService.UpdateTask(&task); err != nil {
		if err.Error() == "task not found" {
			http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{})
}

// Удаление задачи 
func (h *TaskHandler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if err := h.TaskService.DeleteTask(id); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{})
}

// Обработчик, который делает задачу выполненной
func (h *TaskHandler) DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Missing task ID"}`, http.StatusBadRequest)
		return
	}

	now := time.Now()

	err := h.TaskService.TaskDone(id, now)
	if err != nil {
		if err.Error() == "task not found" {
			http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{})
}