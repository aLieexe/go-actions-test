package http

import (
	"jwt-golang/internal/auth"
	"jwt-golang/internal/common"
	"jwt-golang/internal/middlewares"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes(handlers *Handlers, authManager auth.AuthManager, app *common.Application) http.Handler {
	router := chi.NewRouter()

	// fix middleware order
	router.Use(middlewares.RecoverPanic(app))
	router.Use(middlewares.CommonHeaders)
	router.Use(middlewares.RequestLogger(app))
	router.Use(auth.JwtMiddleware(authManager, app))

	//user
	router.Post("/user", handlers.User.PostUser)
	router.Get("/user", handlers.User.GetAllUser)
	router.Delete("/user", handlers.User.DeleteUser)
	router.Put("/user", handlers.User.EditUser)
	router.Get("/user/cookie", handlers.User.CheckCookie)

	//auth
	router.Post("/auth/login", handlers.Auth.UserLogin)
	router.Post("/auth/refresh", handlers.Auth.RefreshToken)
	router.Get("/auth/check", handlers.Auth.CheckJwt)

	router.Group(func(r chi.Router) {
		r.Use(auth.RequireAuthentication(authManager, app))
		router.Get("/auth/me", handlers.Auth.GetCurrentUser)
	})

	router.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		app.WriteJSON(w, 200, common.Envelope{"status": "success"}, nil)
	})

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.ErrorResponse(w, r, http.StatusNotFound, "resource not found")
	})

	// add method not allowed handler
	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		app.ErrorResponse(w, r, http.StatusMethodNotAllowed, "method not allowed")
	})

	return router
}
