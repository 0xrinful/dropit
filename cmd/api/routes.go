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
	router.HandlerFunc("GET", "/files/:token", app.GetFileHandler)
	router.HandlerFunc("DELETE", "/files/:token", app.requireAuth(app.deleteFileHanlder))

	router.HandlerFunc("POST", "/users", app.registerUserHandler)
	router.HandlerFunc("GET", "/users/:id/files", app.getUserFilesHanlder)

	router.HandlerFunc("POST", "/auth/login", app.loginHandler)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
