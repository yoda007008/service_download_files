package service

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

func (t *TaskManager) Run() {
	for i := 0; i < t.workers; i++ {
		t.waitGroup.Add(1)
		go t.Worker()
	}
}

func (t *TaskManager) Stop() {
	close(t.stopCh)
	t.waitGroup.Wait()
}

func (t *TaskManager) Worker() {
	defer t.waitGroup.Done()
	for {
		select {
		case task := <-t.queue:
			t.RunTask(task)
		case <-t.stopCh:
			return
		}
	}
}

func (t *TaskManager) RunTask(task *dto.Task) {
	t.mu.Lock()
	task.Status = dto.StatusRunning
	t.mu.Unlock()

	taskDir := filepath.Join("downloads", task.ID)
	os.MkdirAll(taskDir, 0755)

	for i, url := range task.URLs {
		t.mu.Lock()
		if task.Status == dto.StatusFailed { // отменена
			t.mu.Unlock()
			return
		}
		t.mu.Unlock()

		if err := t.DownloadFile(taskDir, url); err != nil {
			t.mu.Lock()
			task.Status = dto.StatusFailed
			task.Error = err.Error()
			t.mu.Unlock()
			t.Save()
			return
		}

		t.mu.Lock()
		task.Completed = i + 1
		t.mu.Unlock()
		t.Save()
	}

	t.mu.Lock()
	task.Status = dto.StatusDone
	t.mu.Unlock()
	t.Save()
}

func (m *TaskManager) DownloadFile(dir, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	filename := filepath.Join(dir, filepath.Base(url))
	out, err := os.Create(filename)
	if err != nil {
		return err // todo обработка ошибки
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err // todo обработка ошибки
}
