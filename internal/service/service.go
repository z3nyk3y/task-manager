package service

import (
	"github.com/z3nyk3y/task-manager/pkg/workerpool"
)

type Service struct {
	TaskService
}

type Repo struct {
	TaskRepo taskRepo
}

func NewService(repos Repo, workerPool *workerpool.WorkerPool, deadline int) *Service {
	return &Service{
		TaskService: *NewTaskService(repos.TaskRepo, workerPool, deadline),
	}
}
