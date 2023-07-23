package main

import (
	"context"
	"github.com/wisesight/go-api-template/pkg/usecase"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/wisesight/go-api-template/cmd/api/handler"
	"github.com/wisesight/go-api-template/cmd/api/middleware"
	"github.com/wisesight/go-api-template/cmd/api/route"
	"github.com/wisesight/go-api-template/config"
	"github.com/wisesight/go-api-template/pkg/adapter"
	"github.com/wisesight/go-api-template/pkg/log"
	"github.com/wisesight/go-api-template/pkg/repository"
	"github.com/wisesight/go-api-template/pkg/validator"
)

// @title Wisesight API Template
// @version 1.0
// @description This is a sample Wisesight API template server.
// @termsOfService https://wisesight.dev/terms/

// @contact.name API Support
// @contact.url https://wisesight.dev/support
// @contact.email dev@wisesight.com

// @host api.wisesight.dev
// @BasePath /v1

// @schemes https
func main() {
	cfg := config.NewConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongodbClient, err := adapter.NewMongoDBConnection(ctx, cfg.MongoDBURI)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = mongodbClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	if err = validator.NewValidator(); err != nil {
		panic(err)
	}

	userCollection := mongodbClient.Database("test").Collection("users")
	mongoDBAdapter := adapter.NewMongoDBAdapter(mongodbClient)

	userConfig := repository.UserConfig{
		Timeout: 10 * time.Second,
	}
	repository.NewUser(userConfig, mongoDBAdapter, userCollection)

	app := echo.New()

	logger, err := log.NewLoggerZap(&log.ZapConfig{Debug: true})

	if err != nil {
		panic(err)
	}
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLoggerMiddleware(logger))
	app.Use(middleware.ResponseLoggerMiddleware(logger))
	app.Use(middleware.SecurityMiddleware())
	app.Use(middleware.CorsMiddleware())

	timelineRepo := repository.NewTimeline(mongoDBAdapter)

	timelineUseCase := usecase.NewTimeline(timelineRepo)

	userHandler := handler.NewUser(logger)
	timelineHandler := handler.NewTimeline(timelineUseCase, logger)
	probeHandler := handler.NewProbe(mongoDBAdapter, logger)

	route.NewRoute(cfg, app, userHandler, probeHandler, timelineHandler)

	err = app.Start(":5555")
	if err != nil {
		go func() {
			if err := app.Start(":5555"); err != nil && err != http.ErrServerClosed {
				app.Logger.Fatal("shutting down the server")
				panic(err)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		defer cancel()
		if err := app.Shutdown(ctx); err != nil {
			app.Logger.Fatal(err)
		}
	}
}
