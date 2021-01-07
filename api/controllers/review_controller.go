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

//Controller to create new review
func (server *Server) CreateReview(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	review := models.Review{}
	err = json.Unmarshal(body, &review)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	review.Prepare()
	err = review.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != review.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	reviewCreated, err := review.UploadReview(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, reviewCreated.ID))
	responses.JSON(w, http.StatusCreated, reviewCreated)

}

//Controller to get all reviews
func (server *Server) GetReviews(w http.ResponseWriter, r *http.Request) {

	review := models.Review{}

	reviews, err := review.FindAllReviews(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, reviews)
}

//Controller to get specific user's reviews
func (server *Server) GetUserReviews(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	review := models.Review{}

	reviewsRecieved, err := review.FindUserReviews(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, reviewsRecieved)
}

//Controller to update existing review
func (server *Server) UpdateReview(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the review id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//Check if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the review exist
	review := models.Review{}
	err = server.DB.Debug().Model(models.Review{}).Where("id = ?", pid).Take(&review).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Review not found"))
		return
	}

	// If a user attempt to update a review not belonging to him
	if uid != review.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	reviewUpdate := models.Review{}
	err = json.Unmarshal(body, &reviewUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != reviewUpdate.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	reviewUpdate.Prepare()
	err = reviewUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	reviewUpdate.ID = review.ID //this is important to tell the model the trbiew id to update, the other update field are set above

	reviewUpdated, err := reviewUpdate.UpdateReview(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, reviewUpdated)
}

//Controller to delete review
func (server *Server) DeleteReview(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid review id given to us?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Check user authentication
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the review exist
	review := models.Review{}
	err = server.DB.Debug().Model(models.Review{}).Where("id = ?", pid).Take(&review).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Review Not Foundkp"))
		return
	}

	// Is the authenticated user, the owner of this review
	if uid != review.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	_, err = review.DeleteReview(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))

	response := map[string]string{
		"message": "Review deleted",
	}

	responses.JSON(w, http.StatusNoContent, response)
}
