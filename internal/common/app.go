package common

import (
	"log/slog"
)

type Application struct {
	Logger *slog.Logger
	Config *Config
}
