package service

import (
	"context"
	"testing"

	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/z3nyk3y/task-manager/internal/models"
	"github.com/z3nyk3y/task-manager/pkg/workerpool"
	"go.uber.org/zap"
)

func TestProcessTasks(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatal(err)
	}

	zap.ReplaceGlobals(logger)

	ctx := context.TODO()

	taskRepoMock := newMocktaskRepo(t)

	wp, _ := workerpool.New(ctx, 10, 10)

	repos := Repo{
		TaskRepo: taskRepoMock,
	}

	svc := NewService(repos, wp)

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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskRepoMock.On("UpdateTask", mock.Anything, mock.Anything).Maybe().Return(nil)

			err := svc.TaskService.ProcessTasks(ctx, tt.numberOfTasks, tt.processTimeMin, tt.processTimeMax, tt.sucessProbability)

			require.Nil(t, err)

		})
	}
}

var tests = []struct {
	name string

	numberOfTasks     int
	processTimeMin    int64
	processTimeMax    int64
	sucessProbability int
}{
	{
		name: "common test",

		numberOfTasks:     10,
		processTimeMin:    1000,
		processTimeMax:    1001,
		sucessProbability: 100,
	},
}
