package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

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

func (app *application) GetFileHandler(w http.ResponseWriter, r *http.Request) {
	token, err := app.readTokenParam(r)
	if err != nil {
		app.sendNotFoundError(w, r)
		return
	}

	fileInfo, err := app.models.Files.GetByToken(token)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.sendNotFoundError(w, r)
		default:
			app.sendServerError(w, r, err)
		}
		return
	}

	file, err := os.Open(filepath.Join(uploadDir, fileInfo.StoragePath))
	if err != nil {
		app.sendServerError(w, r, err)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		app.sendServerError(w, r, err)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(fileInfo.Filename))
	http.ServeContent(w, r, fileInfo.Filename, stat.ModTime(), file)

	fileInfo.LastAccessedAt = time.Now()
	fileInfo.DownloadCount += 1

	err = app.models.Files.Update(fileInfo)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.sendEditConflictError(w, r)
		default:
			app.sendServerError(w, r, err)
		}
	}
}

func (app *application) deleteFileHanlder(w http.ResponseWriter, r *http.Request) {
	token, err := app.readTokenParam(r)
	if err != nil {
		app.sendNotFoundError(w, r)
		return
	}

	file, err := app.models.Files.GetByToken(token)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.sendNotFoundError(w, r)
		default:
			app.sendServerError(w, r, err)
		}
		return
	}

	user := app.contextGetUser(r)
	if *file.OwnerID != user.ID {
		app.sendPermissionDeniedError(w, r)
		return
	}

	err = app.models.Files.Delete(token, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.sendNotFoundError(w, r)
		default:
			app.sendServerError(w, r, err)
		}
		return
	}

	filepath := filepath.Join(uploadDir, file.StoragePath)
	err = os.Remove(filepath)
	if err != nil {
		app.logger.PrintError(fmt.Errorf(
			"failed to delete file from disk path=%s: %v",
			filepath, err,
		))
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "file successfully deleted"}, nil)
	if err != nil {
		app.sendServerError(w, r, err)
	}
}
