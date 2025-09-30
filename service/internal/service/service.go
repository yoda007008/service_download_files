package service

import (
	"sync"
	"testtask/service/internal/dto"
)

type TaskManager struct {
	mu         sync.Mutex
	tasks      map[string]*dto.Task
	queue      chan *dto.Task
	stopCh     chan struct{}
	filePath   string
	workers    int
	waitGroupe sync.WaitGroup
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
