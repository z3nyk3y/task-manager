package handler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateTaskRequest(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTaskRequest(tt.taskRequest)
			if tt.isError {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}

		})
	}
}

var tests = []struct {
	name string

	taskRequest TaskRequest

	isError bool
}{
	{
		name: "valid test",

		taskRequest: TaskRequest{
			NumberOfTasks:      10,
			ProcessTimeMinimum: 100,
			ProcessTimeMax:     100,
			SucessProbability:  99,
		},

		isError: false,
	},
	{
		name: "invalid NumberOfTasks",

		taskRequest: TaskRequest{
			NumberOfTasks:      -1,
			ProcessTimeMinimum: 100,
			ProcessTimeMax:     100,
			SucessProbability:  99,
		},

		isError: true,
	},
	{
		name: "invalid ProcessTimeMinimum",

		taskRequest: TaskRequest{
			NumberOfTasks:      10,
			ProcessTimeMinimum: -1,
			ProcessTimeMax:     100,
			SucessProbability:  99,
		},

		isError: true,
	},
	{
		name: "invalid ProcessTimeMax",

		taskRequest: TaskRequest{
			NumberOfTasks:      10,
			ProcessTimeMinimum: 100,
			ProcessTimeMax:     -1,
			SucessProbability:  99,
		},

		isError: true,
	},
	{
		name: "invalid SucessProbability",

		taskRequest: TaskRequest{
			NumberOfTasks:      10,
			ProcessTimeMinimum: 100,
			ProcessTimeMax:     -1,
			SucessProbability:  101,
		},

		isError: true,
	},
	{
		name: "ProcessTimeMinimum grater than ProcessTimeMax",

		taskRequest: TaskRequest{
			NumberOfTasks:      10,
			ProcessTimeMinimum: 11,
			ProcessTimeMax:     10,
			SucessProbability:  101,
		},

		isError: true,
	},
}
