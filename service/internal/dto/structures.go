package dto

type TaskStatus string

// todo maybe new status
const (
	StatusPending TaskStatus = "pending"
	StatusRunning TaskStatus = "in_progress"
	StatusDone    TaskStatus = "done"
	StatusFailed  TaskStatus = "failed"
)

type Task struct {
	ID     string     `json:"id"`
	URLs   []string   `json:"urls"`
	Status TaskStatus `json:"status"`
}
