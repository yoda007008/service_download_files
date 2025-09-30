package handlers

import (
	"encoding/json"
	"net/http"
	"testtask/service/internal/service"
)

type TaskHandlers struct {
	manager service.TaskManager
}

func (h *TaskHandlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URLs []string `json:"urls"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	task := h.manager.CreateTask(req.URLs)
	writeJSON(w, task, http.StatusCreated)
}
