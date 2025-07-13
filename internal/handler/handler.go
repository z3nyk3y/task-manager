package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/z3nyk3y/task-manager/internal/service"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	router   *echo.Echo
	services *service.Service
}

type Config struct {
	Host string
	Port string
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services, router: echo.New()}
}

func (h *Handler) ListenAndServe(ctx context.Context, cfg Config) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

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
