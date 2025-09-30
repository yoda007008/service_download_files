package service

import (
	"encoding/json"
	"github.com/google/uuid"
	"os"
	"sync"
	"testtask/service/internal/dto"
	"time"
)

type TaskManager struct {
	mu        sync.Mutex
	tasks     map[string]*dto.Task
	queue     chan *dto.Task
	stopCh    chan struct{}
	filePath  string
	workers   int
	waitGroup sync.WaitGroup
}

func NewTaskManager(filePath string, workers int) *TaskManager {
	return &TaskManager{
		tasks:    make(map[string]*dto.Task),
		queue:    make(chan *dto.Task, 100),
		stopCh:   make(chan struct{}),
		filePath: filePath,
		workers:  workers,
	}
}

func (t *TaskManager) CreateTask(urls []string) *dto.Task {
	t.mu.Lock()
	defer t.mu.Unlock()

	task := &dto.Task{
		ID:        uuid.New().String(),
		URLs:      urls,
		Status:    dto.StatusPending,
		Total:     len(urls),
		CreatedAt: time.Now(),
	}

	t.tasks[task.ID] = task
	t.queue <- task
	t.Save()
	return task
}

func (t *TaskManager) Save() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	f, err := os.Create(t.filePath)
	if err != nil {
		return err // todo обработка ошибки
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(t.tasks)
}
