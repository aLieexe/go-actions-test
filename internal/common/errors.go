package common

import (
	"errors"
	"net/http"
)

var (
	// database errors
	ErrQueryError     = errors.New("query error")
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")

	// authentication errors
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrSessionExpired     = errors.New("session expired")
	ErrUnauthorized       = errors.New("unauthorized access")

	// validation errors
	ErrInvalidInput = errors.New("invalid input")
	ErrMissingField = errors.New("required field missing")
)

func (app *Application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.Logger.Error(err.Error(), "method", method, "URI", uri)
}

// generic return error response
func (app *Application) ErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, message any) {
	envelope := Envelope{"error": message}

	err := app.WriteJSON(w, statusCode, envelope, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}
