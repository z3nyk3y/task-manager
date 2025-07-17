package service

import (
	"context"
	"testing"

	"github.com/z3nyk3y/task-manager/internal/models"
	"github.com/z3nyk3y/task-manager/pkg/workerpool"

	mock "github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestProcessTasks(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	zap.ReplaceGlobals(logger)

	ctx := context.TODO()

	for _, tt := range tests {
		taskRepoMock := newMocktaskRepo(t)

		wp, _ := workerpool.New(ctx, 10, 10)

		repos := Repo{
			TaskRepo: taskRepoMock,
		}

		taskRepoMock.
			On("FetchTasks", mock.Anything, mock.AnythingOfType("int")).Once().
			Return(func(ctx context.Context, numberOfTasks int) ([]models.Task, error) {

				result := make([]models.Task, 0, numberOfTasks)

				for i := range numberOfTasks {
					result = append(result,
						models.Task{
							Id:     int64(i),
							Status: models.New,
						},
					)
				}

				return result, nil
			})

		taskRepoMock.On("UpdateTask", mock.Anything, mock.Anything).Maybe().Return(nil)

		svc := NewService(repos, wp, 2)
		t.Run(tt.name, func(t *testing.T) {
			svc.TaskService.ProcessTasks(ctx, tt.numberOfTasks, tt.processTimeMin, tt.processTimeMax, tt.sucessProbability)
		})
	}
}

var tests = []struct {
	name string

	numberOfTasks     int
	processTimeMin    int
	processTimeMax    int
	sucessProbability int

	isError bool
}{
	{
		name: "common test",

		numberOfTasks:     10,
		processTimeMin:    1000,
		processTimeMax:    1001,
		sucessProbability: 100,

		isError: false,
	},
	{
		name: "common test",

		numberOfTasks:     20,
		processTimeMin:    1000,
		processTimeMax:    1000,
		sucessProbability: 100,

		isError: true,
	},
	{
		name: "more than deadline",

		numberOfTasks:     1,
		processTimeMin:    100000,
		processTimeMax:    100000,
		sucessProbability: 100,

		isError: true,
	},
	{
		name: "no time slee",

		numberOfTasks:     10,
		processTimeMin:    0,
		processTimeMax:    0,
		sucessProbability: 100,

		isError: true,
	},
}
