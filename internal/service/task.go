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
	deadline   int
}

func NewTaskService(repo taskRepo, workerPool *workerpool.WorkerPool, deadline int) *TaskService {
	return &TaskService{repo: repo, workerPool: workerPool, deadline: deadline}
}

func (ts *TaskService) ProcessTasks(ctx context.Context, numberOfTasks, processTimeMinimum, processTimeMax int, sucessProbability int) {
	logger := zap.L()

	tasks, err := ts.repo.FetchTasks(ctx, numberOfTasks)
	if err != nil {
		logger.Info("unable to fetch tasks", zap.Error(err))
		return
	}

	logger.Info(fmt.Sprintf("proccess tasks %v", tasks))

	resultCh := make(chan models.Task, numberOfTasks)

	var wg sync.WaitGroup

	for _, task := range tasks {
		task := task

		deadline := time.NewTicker(time.Duration(ts.deadline) * time.Second)

		p := float64(sucessProbability) / 100.0

		var taskStatus models.TaskStatus

		processTimeMinimum := int64(processTimeMinimum)
		processTimeMax := int64(processTimeMax)
		wg.Add(1)
		err := ts.workerPool.AddJob(func() {
			defer wg.Done()

			if rand.Float64() < p {
				taskStatus = models.Processed
			} else {
				taskStatus = models.New
			}

			delta := processTimeMax - processTimeMinimum

			var taskTime int64

			if delta == 0 {
				taskTime = processTimeMinimum
			} else {
				taskTime = processTimeMinimum + rand.Int63n(delta)
			}

			for taskTime >= 0 {
				select {
				case <-deadline.C:
					logger.Info(fmt.Sprintf("deadline excedeed for task %d, setting status to %s", task.Id, models.New))
					task.Status = models.New
					return
				default:
					if taskTime < 10 {
						time.Sleep(time.Duration(time.Duration(taskTime) * time.Millisecond))
					} else {
						time.Sleep(time.Duration(10 * time.Millisecond))
					}
					taskTime -= 10
				}
			}

			task.Status = taskStatus
			resultCh <- task
		})
		if err != nil {
			wg.Done()
			task.Status = models.New
			err := ts.repo.UpdateTask(ctx, task)
			if err != nil {
				logger.Error("error restore task status", zap.Int64("Id", task.Id), zap.String("Status", string(task.Status)))
			} else {
				logger.Info(fmt.Sprintf("chan is full; restored task %d with status %s", task.Id, task.Status))
			}
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

}

func randTime(start, end time.Time) time.Time {
	diff := end.Sub(start)

	rand := time.Duration(rand.Int63n(int64(diff)))

	return start.Add(rand)
}
