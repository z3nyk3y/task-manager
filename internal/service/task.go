package service

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/z3nyk3y/task-manager/internal/models"
	"github.com/z3nyk3y/task-manager/pkg/workerpool"
	"go.uber.org/zap"
)

type taskRepo interface {
	FetchTasks(ctx context.Context, numberOfTasks int) ([]models.Task, error)
	UpdateTask(ctx context.Context, tasks models.Task) error
}

type TaskService struct {
	repo       taskRepo
	workerPool *workerpool.WorkerPool
}

func NewTaskService(repo taskRepo, workerPool *workerpool.WorkerPool) *TaskService {
	return &TaskService{repo: repo, workerPool: workerPool}
}

func (ts *TaskService) ProcessTasks(ctx context.Context, numberOfTasks int, processTimeMin, processTimeMax int64, sucessProbability int) error {
	logger := zap.L()

	tasks, err := ts.repo.FetchTasks(ctx, numberOfTasks)
	if err != nil {
		return err
	}

	var test bool
	logger.Info(fmt.Sprintf("test %v", test))

	logger.Info(fmt.Sprintf("proccess tasks %v", tasks))

	resultCh := make(chan models.Task, numberOfTasks)

	var wg sync.WaitGroup
	wg.Add(numberOfTasks)

	for _, task := range tasks {
		// task := task

		p := float64(sucessProbability) / 100.0

		var taskStatus models.TaskStatus

		processTimeMin := processTimeMin
		processTimeMax := processTimeMax

		err := ts.workerPool.AddJob(func() {
			defer wg.Done()

			test = true

			if rand.Float64() < p {
				taskStatus = models.Processed
			} else {
				taskStatus = models.New
			}

			delta := processTimeMax - processTimeMin

			time.Sleep(time.Duration(processTimeMin+rand.Int63n(delta)) * time.Millisecond)

			task.Status = taskStatus
			resultCh <- task
		})
		if err != nil {
			return err
		}
	}

	go func() {
		wg.Wait()

		close(resultCh)
	}()

	for task := range resultCh {
		err := ts.repo.UpdateTask(ctx, task)
		if err != nil {
			logger.Error("error updating task", zap.Int64("Id", task.Id), zap.String("Status", string(task.Status)))
		} else {
			logger.Info(fmt.Sprintf("updated task %d with status %s", task.Id, task.Status))
		}
	}

	return nil
}

func randTime(start, end time.Time) time.Time {
	diff := end.Sub(start)

	rand := time.Duration(rand.Int63n(int64(diff)))

	return start.Add(rand)
}
