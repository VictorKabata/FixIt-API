package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/victorkabata/FixIt-API/api/auth"
	"github.com/victorkabata/FixIt-API/api/models"
	"github.com/victorkabata/FixIt-API/api/responses"
	"github.com/victorkabata/FixIt-API/api/utils/formaterror"
	// "github.com/victorkabata/golang-authentication/api/auth"
	// "github.com/victorkabata/golang-authentication/api/responses"
	// "github.com/victorkabata/golang-authentication/api/utils/formaterror"
	// "github.com/victorkabata/golang-authentication/models"
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
