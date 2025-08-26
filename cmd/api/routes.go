package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.sendNotFoundError)
	router.MethodNotAllowed = http.HandlerFunc(app.sendMethodNotAllowedError)

	router.HandlerFunc("POST", "/files", nil)
	router.HandlerFunc("GET", "/files/:token", nil)
	router.HandlerFunc("DELETE", "/files/:token", nil)

	router.HandlerFunc("POST", "/users", app.registerUserHandler)
	router.HandlerFunc("GET", "/users/:id/files", nil)

	router.HandlerFunc("POST", "/auth/login", nil)

	return router
}
