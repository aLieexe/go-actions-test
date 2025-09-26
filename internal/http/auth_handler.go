package http

import (
	"jwt-golang/internal/auth"
	"jwt-golang/internal/common"
	"jwt-golang/internal/models"
	"net/http"
)

type AuthHandlerInterface interface {
	UserLogin(w http.ResponseWriter, r *http.Request)
	CheckJwt(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
	GetCurrentUser(w http.ResponseWriter, r *http.Request)
}

type AuthHandler struct {
	userModel   models.UserModelInterface
	authManager auth.AuthManager
	app         *common.Application
}

func NewAuthHandler(userModel models.UserModelInterface, authManager auth.AuthManager, app *common.Application) *AuthHandler {
	return &AuthHandler{
		userModel:   userModel,
		authManager: authManager,
		app:         app,
	}
}

func (h *AuthHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	userInput := models.UserLogin{}
	err := h.app.ReadJSON(w, r, &userInput)
	if err != nil {
		h.app.ErrorResponse(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user, err := h.userModel.AuthenticateUser(r.Context(), userInput)
	if err != nil {
		h.app.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := h.authManager.GenerateTokens(r.Context(), user)
	if err != nil {
		h.app.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = h.app.WriteJSON(w, http.StatusOK, common.Envelope{"access_token": tokens[0], "refresh_token": tokens[1]}, nil)
	if err != nil {
		h.app.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *AuthHandler) CheckJwt(w http.ResponseWriter, r *http.Request) {
	claims, exist := auth.GetUserFromContext(r.Context())
	if !exist {
		h.app.ErrorResponse(w, r, http.StatusUnauthorized, "jwt doesnt exist")
		return
	}

	err := h.app.WriteJSON(w, http.StatusOK, common.Envelope{"data": claims}, nil)
	if err != nil {
		h.app.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}

	err := h.app.ReadJSON(w, r, &input)
	if err != nil {
		h.app.Logger.Info(err.Error())

		h.app.ErrorResponse(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}
	h.app.Logger.Info(input.RefreshToken)

	userId, err := h.authManager.GetUserIdFromRefreshToken(r.Context(), input.RefreshToken)
	if err != nil {
		h.app.ErrorResponse(w, r, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := h.userModel.GetUserById(r.Context(), userId)
	if err != nil {
		h.app.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	tokens, err := h.authManager.GenerateTokens(r.Context(), user)
	if err != nil {
		h.app.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = h.app.WriteJSON(w, http.StatusOK, common.Envelope{"access_token": tokens[0], "refresh_token": tokens[1]}, nil)
	if err != nil {
		h.app.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userId, exist := auth.GetUserIdFromContext(r.Context())
	if !exist {
		h.app.ErrorResponse(w, r, http.StatusUnauthorized, "jwt doesnt exist")
		return
	}

	user, err := h.userModel.GetUserById(r.Context(), userId)
	if err != nil {
		h.app.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
	}

	err = h.app.WriteJSON(w, http.StatusOK, common.Envelope{"user": user}, nil)
	if err != nil {
		h.app.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

}
