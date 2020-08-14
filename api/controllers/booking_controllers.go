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

func (server *Server) MakeBooking(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	booking := models.Booking{}
	err = json.Unmarshal(body, &booking)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	booking.Prepare()
	err = booking.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if uid != booking.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	bookingMade, err := booking.SaveBooking(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, bookingMade.ID))
	responses.JSON(w, http.StatusCreated, bookingMade)
}

func (server *Server) GetBookings(w http.ResponseWriter, r *http.Request) {

	booking := models.Booking{}

	bookings, err := booking.FindAllBookings(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, bookings)
}

func (server *Server) UpdateBooking(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the booking id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Check if the post exist
	booking := models.Booking{}
	err = server.DB.Debug().Model(models.Booking{}).Where("id = ?", pid).Take(&booking).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Booking not found"))
		return
	}

	// Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	bookingUpdate := models.Booking{}

	err = json.Unmarshal(body, &bookingUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	bookingUpdate.ID = booking.ID //this is important to tell the model the post id to update, the other update field are set above

	bookingUpdated, err := bookingUpdate.UpdateABooking(server.DB, pid)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, bookingUpdated)
}
