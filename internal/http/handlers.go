package http

import (
	"jwt-golang/internal/auth"
	"jwt-golang/internal/common"
	"jwt-golang/internal/models"
)

type Handlers struct {
	User UserHandlerInterface
	Auth AuthHandlerInterface
}

func NewHandlers(models models.Models, authManager auth.AuthManager, app *common.Application) Handlers {
	return Handlers{
		User: NewUserHandler(models.Users, app),
		Auth: NewAuthHandler(models.Users, authManager, app),
	}
}
