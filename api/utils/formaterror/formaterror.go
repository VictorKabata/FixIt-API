package formaterror

import (
	"errors"
	"strings"
)

//Formats error to a more readable format
func FormatError(err string) error {

	if strings.Contains(err, "email") {
		return errors.New("Email Already Taken")
	}

	if strings.Contains(err, "hashedPassword") {
		return errors.New("Incorrect Password")
	}
	return errors.New("Incorrect Details")
}
