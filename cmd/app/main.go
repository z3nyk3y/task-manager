package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/z3nyk3y/task-manager/internal/handler"
	"github.com/z3nyk3y/task-manager/internal/repo/postgresql"
	"github.com/z3nyk3y/task-manager/internal/service"
	"github.com/z3nyk3y/task-manager/pkg/workerpool"

	"go.uber.org/zap"
)

const (
	apiPort    = "API_PORT"
	apiHost    = "API_HOST"
	Login      = "DB_LOGIN"
	DbName     = "DB_NAME"
	DbHost     = "DB_HOST"
	DbPort     = "DB_PORT"
	DbPassword = "DB_PASSWORD"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := StartTaskManager(ctx)
	if err != nil {
		log.Fatalf("error while sturting up startTaskManager: %s", err.Error())
	}
}

func StartTaskManager(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("logger creation error %s", err.Error())
	}
	zap.ReplaceGlobals(logger)

	repoCfg := postgresql.Config{
		Login:    os.Getenv("DB_LOGIN"),
		Name:     os.Getenv("DB_NAME"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	pool, err := postgresql.NewPostgreSqlDB(ctx, repoCfg)
	if err != nil {
		log.Fatalf("unable to estaiblish connect with database %s", err.Error())
	}
	defer pool.Close()

	repo := postgresql.NewRepository(pool)

	wp, err := workerpool.New(ctx, 100, 100)
	if err != nil {
		log.Fatalf("unable to create worker pool %s", err.Error())
	}

	services := service.NewService(
		service.Repo{
			TaskRepo: repo.Task,
		},
		wp,
		60,
	)

	handlers := handler.NewHandler(services, wp)

	handlersCfg := handler.Config{
		Host: os.Getenv(apiHost),
		Port: os.Getenv(apiPort),
	}

	go func() {
		defer cancel()

		err := handlers.ListenAndServe(handlersCfg)
		if err != nil {
			logger.Error("error starting up server", zap.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down server in progress")

	err = handlers.ShutDown()
	if err != nil {
		logger.Error("could not shut down server")
	}

	err = logger.Sync()
	if err != nil {
		return err
	}

	return nil
}
