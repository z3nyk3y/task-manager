package main

import (
	"context"
	"log"
	"os"

	"github.com/z3nyk3y/task-manager/internal/handler"
	"github.com/z3nyk3y/task-manager/internal/repo/postgresql"
	"github.com/z3nyk3y/task-manager/internal/service"
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

	services := service.NewService(service.Repo{
		TaskRepo: repo.Task,
	})

	handlers := handler.NewHandler(services)

	handlersCfg := handler.Config{
		Host: os.Getenv(apiHost),
		Port: os.Getenv(apiPort),
	}

	go func() {
		defer cancel()

		err := handlers.ListenAndServe(ctx, handlersCfg)
		if err != nil {
			logger.Error("error starting up server", zap.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down server in porgress")

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
