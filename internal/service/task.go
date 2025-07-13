package service

import (
	"context"

	"gitlab.com/eightlix-group/task-manager/internal/models"
)

type taskRepo interface {
	GetTasks(ctx context.Context) ([]models.Task, error)
}

type Task struct {
	repo taskRepo
}

func NewTaskService(repo taskRepo) *Task {
	return &Task{repo: repo}
}
