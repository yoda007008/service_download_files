package handlers

import (
	"encoding/json"
	"net/http"
	"testtask/service/internal/service"
)

func NewTaskHandlers(manager *service.TaskManager) *TaskHandlers {
	return &TaskHandlers{manager: *manager}
}

func (h *TaskHandlers) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/tasks", h.CreateTask)

}

func writeJSON(w http.ResponseWriter, data any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
