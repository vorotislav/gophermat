package main

import (
	"context"
	"go.uber.org/zap"
	"gophermat/internal/app"
	"gophermat/internal/authentication"
	"gophermat/internal/http"
	"gophermat/internal/http/client"
	"gophermat/internal/repository/postgres"
	"gophermat/internal/settings"
	"gophermat/internal/signals"
	"log"
	"os"
	"time"
)

const serviceShutdownTimeout = 1 * time.Second

func main() {
	set := settings.Settings{}

	parseFlag(&set)

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Printf("cannot create logger: %s", err.Error())

		os.Exit(1)
	}

	defer logger.Sync()

	logger.Debug("Server starting...")
	logger.Debug("Current settings",
		zap.String("ip address", set.Address),
		zap.String("database uri", set.DatabaseURI),
		zap.String("accrual system address", set.AccrualSystemAddress))

	ctx, cancel := context.WithCancel(context.Background())
	oss := signals.NewOSSignals(ctx)

	oss.Subscribe(func(sig os.Signal) {
		logger.Info("Stopping by OS Signal...",
			zap.String("signal", sig.String()))

		cancel()
	})

	repo, err := postgres.NewStorage(ctx, logger, &set)
	if err != nil {
		logger.Error("create storage", zap.Error(err))

		return
	}

	auth := authentication.NewAuthenticator()

	accrualClient := client.NewClient(logger, &set)

	gm := app.NewGMart(logger, auth, repo, accrualClient)

	s, err := http.NewService(logger, &set, gm, auth)
	if err != nil {
		logger.Error("create http service", zap.Error(err))

		return
	}

	serviceErrCh := make(chan error, 1)
	go func(errCh chan<- error) {
		defer close(errCh)

		if err := s.Run(); err != nil {
			errCh <- err
		}
	}(serviceErrCh)

	select {
	case err := <-serviceErrCh:
		if err != nil {
			logger.Error("service error", zap.Error(err))
			cancel()
		}
	case <-ctx.Done():
		logger.Info("Server stopping...")
		ctxShutdown, ctxCancelShutdown := context.WithTimeout(context.Background(), serviceShutdownTimeout)

		if err := s.Stop(ctxShutdown); err != nil {
			logger.Error("cannot stop server", zap.Error(err))
		}

		repo.Stop()

		defer ctxCancelShutdown()
	}
}
