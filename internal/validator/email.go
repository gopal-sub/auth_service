package validator

import (
	"net/mail"
)


func Email(email string) error{
	if email == ""{
		return EmailRequiredErr
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return InvalidEmailErr
	}

	return nil


}