package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"ozon-fintech/internal/config"
	"ozon-fintech/internal/routes"
	service2 "ozon-fintech/internal/service"
	"ozon-fintech/internal/storage"
	"ozon-fintech/internal/storage/inMemory"
	"ozon-fintech/internal/storage/postgres"
	"syscall"
)

func main() {
	cfg, err := config.LoadConfig("./configs/")
	if err != nil {
		log.Fatalln(err)
	}
	var dbFlag bool
	flag.BoolVar(&dbFlag, "db", false, "Run with DB postgres: ")
	flag.Parse()

	//TODO: logger
	logger := setupLogger(cfg.Env)

	//TODO: storage
	var dbType storage.DbType
	logger.Info("creating storage")
	if dbFlag {
		db, err := postgres.New(&cfg.PgStorage)
		if err != nil {
			logger.Error("creating new pg storage err: ", err)
			panic(err)
		}

		dbType = db

	} else {
		mem := inMemory.New()
		dbType = mem

	}
	store := storage.New(dbType, logger)
	service := service2.NewService(dbType)

	//TODO: router
	errChan := make(chan error, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	r := routes.New(cfg.ServerPort, cfg.Env, logger, service)

	logger.Info("Starting the server on address:", cfg.ServerPort)

	go func() {
		if err := r.Start(); err != nil {
			errChan <- err
		}
	}()

	select {
	case sig := <-sigChan:
		logger.Info("got interrtupt", sig)
	case err := <-errChan:
		logger.Info("got an error", err)

	}
	logger.Info("Shutdown server")

	ctx := context.Background()

	if err := r.Shutdown(ctx); err != nil {
		logger.Info("Shutdown server err: ", err)
	}

	logger.Info("closing db conn")
	if err := store.Stop(); err != nil {
		logger.Info("closing db conn err: ", err)
	}

	logger.Info("db conn closed")

}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case "local":
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case "prod":
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return logger
}
