package responses

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/victorkabata/FixIt-API/api/auth"
	"github.com/victorkabata/FixIt-API/api/models"
)

func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

func ERROR(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		JSON(w, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	JSON(w, http.StatusBadRequest, nil)
}

func PrepareResponse(user *models.User) map[string]interface{} {
	responseUser := &models.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
	}

	token, _ := auth.CreateToken(user.ID)
	var response = map[string]interface{}{"message": "Successful"}
	response["jwt"] = token
	response["user"] = responseUser

	return response
}
