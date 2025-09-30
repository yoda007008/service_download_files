package dto

import "time"

type TaskStatus string

// todo maybe new status
const (
	StatusPending TaskStatus = "pending"
	StatusRunning TaskStatus = "in_progress"
	StatusDone    TaskStatus = "done"
	StatusFailed  TaskStatus = "failed"
)

type Task struct {
	ID        string     `json:"id"`
	URLs      []string   `json:"urls"`
	Status    TaskStatus `json:"status"`
	Total     int        `json:"total"`
	Completed int        `json:"completed"`
	CreatedAt time.Time  `json:"created_at"`
	Error     string     `json:"error"`
}

type TaskImpl interface { // todo dependency for interface for testing

}
