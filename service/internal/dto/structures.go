package dto

type TaskStatus string

type Task struct {
	ID     string     `json:"id"`
	URLs   []string   `json:"urls"`
	Status TaskStatus `json:"status"`
}
