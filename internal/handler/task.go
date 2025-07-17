package handler

import (
	"context"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
)

type TaskRequest struct {
	NumberOfTasks      int `json:"number_of_tasks"`
	ProcessTimeMinimum int `json:"process_time_minimum"`
	ProcessTimeMax     int `json:"process_time_max"`
	SucessProbability  int `json:"sucess_probability"`
}

func (h *Handler) TaskHandler(c echo.Context) error {
	ctx := context.Background()
	logger := zap.L()

	var taskRequest TaskRequest
	err := c.Bind(&taskRequest)
	if err != nil {
		logger.Error("error bunding data", zap.Error(err))
		return h.sendResponse(c, http.StatusBadRequest, "", "invalid json")
	}

	err = validateTaskRequest(taskRequest)
	if err != nil {
		return h.sendResponse(c, http.StatusBadRequest, "", err.Error())
	}

	err = h.workerPool.AddJob(func() {
		h.services.TaskService.ProcessTasks(ctx, taskRequest.NumberOfTasks, taskRequest.ProcessTimeMinimum, taskRequest.ProcessTimeMax, taskRequest.SucessProbability)
	})
	if err != nil {
		return h.sendResponse(c, http.StatusBadGateway, "", "server is busy")
	}

	return h.sendResponse(c, http.StatusNoContent, "", "")
}

func validateTaskRequest(taskRequest TaskRequest) error {
	if taskRequest.NumberOfTasks <= 0 {
		return errors.New("number_of_tasks must be grater than 0")
	}

	if taskRequest.ProcessTimeMinimum <= 0 {
		return errors.New("process_time_minmum must be grater than 0")
	}
	if taskRequest.ProcessTimeMax <= 0 {
		return errors.New("process_time_max must be grater than 0")
	}

	if taskRequest.ProcessTimeMinimum > taskRequest.ProcessTimeMax {
		return errors.New("process_time_minimum must be less than process_time_max or they must be equal")
	}

	if taskRequest.SucessProbability < 0 || taskRequest.SucessProbability > 100 {
		return errors.New("sucess_probability must be in range from  0 to 100")
	}

	return nil
}
