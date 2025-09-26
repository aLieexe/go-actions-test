package http

import (
	"errors"
	"jwt-golang/internal/auth"
	"jwt-golang/internal/common"
	"jwt-golang/internal/models"
	"net/http"
)

type UserHandlerInterface interface {
	PostUser(w http.ResponseWriter, r *http.Request)
	GetAllUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	EditUser(w http.ResponseWriter, r *http.Request)
	CheckCookie(w http.ResponseWriter, r *http.Request)
}

type UserHandler struct {
	userModel models.UserModelInterface
	app       *common.Application
}

func NewUserHandler(userModel models.UserModelInterface, app *common.Application) *UserHandler {
	return &UserHandler{
		userModel: userModel,
		app:       app,
	}
}

func validateUserInput(user models.PostUser) error {
	if len(user.Username) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if len(user.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
}

func (h *UserHandler) PostUser(w http.ResponseWriter, r *http.Request) {
	userInput := models.PostUser{}
	err := h.app.ReadJSON(w, r, &userInput)
	if err != nil {
		h.app.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = validateUserInput(userInput)
	if err != nil {
		h.app.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := h.userModel.Insert(r.Context(), userInput)
	if err != nil {
		h.app.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = h.app.WriteJSON(w, http.StatusOK, common.Envelope{"user_id": userId}, nil)
	if err != nil {
		h.app.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

}

func (h *UserHandler) GetAllUser(w http.ResponseWriter, r *http.Request) {
	users, err := h.userModel.GetAll(r.Context())
	if err != nil {
		h.app.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = h.app.WriteJSON(w, http.StatusOK, common.Envelope{"data": users}, nil)
	if err != nil {
		h.app.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

}

func (h *UserHandler) EditUser(w http.ResponseWriter, r *http.Request) {
	userInput := models.EditUser{}
	err := h.app.ReadJSON(w, r, &userInput)
	if err != nil {
		h.app.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	currentUserId, ok := auth.GetUserIdFromContext(r.Context())
	if !ok {
		h.app.ErrorResponse(w, r, http.StatusUnauthorized, "user not authenticated")
		return
	}

	err = h.userModel.EditUser(r.Context(), userInput, currentUserId)
	if err != nil {
		h.app.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = h.app.WriteJSON(w, http.StatusOK, common.Envelope{"status": "success"}, nil)
	if err != nil {
		h.app.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}
	err := h.app.ReadJSON(w, r, &input)
	if err != nil {
		h.app.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = h.userModel.DeleteUser(r.Context(), input.Email)
	if err != nil {
		h.app.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = h.app.WriteJSON(w, http.StatusOK, common.Envelope{"status": "success"}, nil)
	if err != nil {
		h.app.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

}

func (h *UserHandler) CheckCookie(w http.ResponseWriter, r *http.Request) {
	envelope := common.Envelope{"data": r.Cookies()}
	err := h.app.WriteJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		h.app.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

}
