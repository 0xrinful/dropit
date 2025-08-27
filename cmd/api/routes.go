package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.sendNotFoundError)
	router.MethodNotAllowed = http.HandlerFunc(app.sendMethodNotAllowedError)

	router.HandlerFunc("POST", "/files", app.uploadHandler)
	router.HandlerFunc("GET", "/files/:token", nil)
	router.HandlerFunc("DELETE", "/files/:token", app.requireAuth(nil))

	router.HandlerFunc("POST", "/users", app.registerUserHandler)
	router.HandlerFunc("GET", "/users/:id/files", nil)

	router.HandlerFunc("POST", "/auth/login", app.loginHandler)

	return app.recoverPanic(app.authenticate(router))
}
