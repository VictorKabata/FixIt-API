package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/victorkabata/FixIt-API/api/auth"
	"github.com/victorkabata/FixIt-API/api/models"
	"github.com/victorkabata/FixIt-API/api/responses"
	"github.com/victorkabata/FixIt-API/api/utils/formaterror"
)

func (server *Server) CreateWork(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	work := models.Work{}
	err = json.Unmarshal(body, &work)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	work.Prepare()
	err = work.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if uid != work.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	workCreated, err := work.UploadWork(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, workCreated.ID))
	responses.JSON(w, http.StatusCreated, workCreated)
}

func (server *Server) GetWork(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	work := models.Work{}

	workReceived, err := work.FindWorkByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, workReceived)
}

func (server *Server) UpdateWork(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the post id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Check if the work exist
	work := models.Work{}
	err = server.DB.Debug().Model(models.Work{}).Where("id = ?", pid).Take(&work).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Work not found"))
		return
	}

	// Read the data posted in json body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	workUpdate := models.Work{}
	err = json.Unmarshal(body, &workUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	workUpdate.ID = work.ID //this is important to tell the model the post id to update, the other update field are set above

	workUpdated, err := workUpdate.UpdateWork(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, workUpdated)
}

//Controller to get user's works
func (server *Server) GetUserWorks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	work := models.Work{}

	userWorks, err := work.FindUserWorks(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, userWorks)
}
