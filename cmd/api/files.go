package main

import (
	"fmt"
	"net/http"

	"github.com/0xrinful/dropit/internal/data"
)

func (app *application) uploadHandler(w http.ResponseWriter, r *http.Request) {
	token, err := data.NewFileToken()
	if err != nil {
		app.sendServerError(w, r, err)
		return
	}

	filename, err := app.uploadFile(w, r, "file", token)
	if err != nil {
		app.sendBadRequestError(w, r, err)
		return
	}

	file := &data.File{
		Token:       token,
		StoragePath: token,
		Filename:    filename,
	}

	user := app.contextGetUser(r)
	if !user.IsAnonymous() {
		file.OwnerID = new(int64)
		*file.OwnerID = user.ID
	}

	err = app.models.Files.Insert(file)
	if err != nil {
		app.sendServerError(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/files/%s", file.Token))

	err = app.writeJSON(w, http.StatusCreated, envelope{"file": file}, headers)
	if err != nil {
		app.sendServerError(w, r, err)
	}
}
