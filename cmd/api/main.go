package main

import (
	"jwt-golang/internal/auth"
	"jwt-golang/internal/common"
	internalHttp "jwt-golang/internal/http"
	"jwt-golang/internal/models"
	"jwt-golang/internal/services"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	config, err := common.Load(".env")

	if err != nil {
		logger.Error(err.Error())
	}
	pg, err := services.ConnectPostgres(config.DatabaseURL)
	if err != nil {
		logger.Error(err.Error())
	}

	newModels := models.NewModels(pg)
	app := &common.Application{
		Logger: logger,
		Config: config,
	}

	jwtProvider := auth.NewJwtProvider(config.AccessTokenSecret, config.RefreshTokenSecret, time.Duration(config.AccesssTokenTTL), time.Duration(config.RefreshTokenTTL))
	authManager := auth.NewAuthManager(jwtProvider)

	newHandlers := internalHttp.NewHandlers(newModels, *authManager, app)

	server := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      internalHttp.Routes(&newHandlers, *authManager, app),
		IdleTimeout:  time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	err = server.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
	}

	os.Exit(1)
}
