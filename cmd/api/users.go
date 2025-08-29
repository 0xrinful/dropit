package main

import (
	"errors"
	"net/http"

	"github.com/0xrinful/dropit/internal/data"
	"github.com/0xrinful/dropit/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.sendBadRequestError(w, r, err)
		return
	}

	user := &data.User{
		Name:  input.Name,
		Email: input.Email,
		Role:  "user",
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.sendServerError(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.sendValidationError(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.sendValidationError(w, r, v.Errors)
		default:
			app.sendServerError(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.sendServerError(w, r, err)
	}
}

func (app *application) getUserFilesHanlder(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.sendNotFoundError(w, r)
		return
	}

	files, err := app.models.Files.GetAllForUser(id)
	if err != nil {
		app.sendServerError(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"files": files}, nil)
	if err != nil {
		app.sendServerError(w, r, err)
	}
}
