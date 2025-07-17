package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/z3nyk3y/task-manager/internal/service"
	"github.com/z3nyk3y/task-manager/pkg/workerpool"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handler struct {
	router     *echo.Echo
	services   *service.Service
	workerPool *workerpool.WorkerPool
}

type Config struct {
	Host string
	Port string
}

type Response struct {
	Status string `json:"status"`
	Result any    `json:"result"`
	Error  string `json:"error"`
}

func NewHandler(services *service.Service, workerPool *workerpool.WorkerPool) *Handler {
	return &Handler{services: services, router: echo.New(), workerPool: workerPool}
}

func (h *Handler) ListenAndServe(cfg Config) error {

	h.router.Use(middleware.Secure())
	h.router.Use(middleware.Recover())

	apiV1 := h.router.Group(("/api/v1"), middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(10)), LogRequestBodyMiddleware)
	{
		apiV1.POST("/tasks", h.TaskHandler)
	}

	err := h.router.Start(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port))
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) ShutDown() error {
	shutDownCtx, shutDownCancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer shutDownCancel()

	err := h.router.Shutdown(shutDownCtx)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) sendResponse(c echo.Context, status int, result any, clientMessage string) error {
	response := Response{
		Result: result,
		Error:  clientMessage,
	}

	return c.JSON(status, response)
}
